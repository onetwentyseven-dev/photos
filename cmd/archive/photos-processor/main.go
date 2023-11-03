package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // Import the image/jpeg package for JPEG support
	_ "image/png"  // Import the image/png package for PNG support
	"math"
	"photos"
	"photos/internal/store/mysql"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rwcarlsen/goexif/exif"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

type EventDetail struct {
	Bucket          *Bucket `json:"bucket"`
	Object          *Object `json:"object"`
	Reason          string  `json:"reason"`
	RequestID       string  `json:"request-id"`
	Requester       string  `json:"requester"`
	SourceIPAddress string  `json:"source-ip-address"`
	Version         string  `json:"version"`
}
type Bucket struct {
	Name string `json:"name"`
}
type Object struct {
	Etag      string `json:"etag"`
	Key       string `json:"key"`
	Sequencer string `json:"sequencer"`
	Size      int    `json:"size"`
}

type app struct {
	logger *logrus.Logger
	s3     *s3.Client
	image  *mysql.ImageRepository
}

func (a *app) handle(ctx context.Context, event *events.CloudWatchEvent) error {

	if event.Version != "0" {
		a.logger.Error("unsupported event version")
		return nil
	}

	if event.Source != "aws.s3" {
		a.logger.Error("unsupported event source")
		return nil
	}

	var detail = new(EventDetail)
	err := json.Unmarshal(event.Detail, detail)
	if err != nil {
		return err
	}

	imageID := strings.Split(detail.Object.Key, ".")[0]

	entry := a.logger.WithFields(logrus.Fields{
		"bucket": detail.Bucket.Name,
		"key":    detail.Object.Key,
	})

	imageMeta, err := a.image.Image(ctx, imageID)
	if err != nil {
		entry.WithError(err).Error("failed to get image details from database")
		return err
	}

	processingFailedFunc := a.processingFailed(ctx, imageMeta, entry)

	var originalBuffer = new(bytes.Buffer)

	objectOutput, err := a.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &detail.Bucket.Name,
		Key:    &detail.Object.Key,
	})
	if err != nil {
		processingFailedFunc(err)
		entry.WithError(err).Error("failed to get object")
		return err
	}

	_, _ = originalBuffer.ReadFrom(objectOutput.Body)

	err = objectOutput.Body.Close()
	if err != nil {
		entry.WithError(err).Error("failed to close object body")
	}

	err = a.processImage(ctx, objectOutput, detail, entry, bytes.NewBuffer(originalBuffer.Bytes()))
	if err != nil {
		processingFailedFunc(err)
		entry.WithError(err).Error("failed to process image")
		return err
	}

	exifData, err := a.extractExifData(ctx, entry, bytes.NewBuffer(originalBuffer.Bytes()))
	if err != nil {
		entry.WithError(err).Error("failed to extract exif data")
	}

	if exifData != nil {
		var o = make(map[string]any)
		err = json.Unmarshal(exifData, &o)
		if err != nil {
			entry.WithError(err).Error("failed to unmarshal exif data")
		}

		imageMeta.ExifData = o
	}

	imageMeta.Status = photos.ProcessedImageStatus

	err = a.image.UpdateImageByImageID(ctx, imageMeta)
	if err != nil {
		processingFailedFunc(err)
		entry.WithError(err).Error("failed to update image status")
		return err
	}

	return nil
}

func (a *app) processingFailed(ctx context.Context, image *photos.Image, entry *logrus.Entry) func(error) {
	return func(err error) {
		if err != nil {
			return
		}

		image.Status = photos.ErroredImageStatus
		image.ProcessingErrors = []string{
			err.Error(),
		}
		err = a.image.UpdateImageByImageID(ctx, image)
		if err != nil {
			entry.WithError(err).Error("failed to update image status")
		}
	}
}

func (a *app) extractExifData(ctx context.Context, entry *logrus.Entry, buf *bytes.Buffer) (json.RawMessage, error) {

	x, err := exif.Decode(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		entry.WithError(err).Error("failed to decode exif data")
		return nil, fmt.Errorf("failed to decode exif data")
	}

	return x.MarshalJSON()

}

func (a *app) processImage(ctx context.Context, object *s3.GetObjectOutput, detail *EventDetail, entry *logrus.Entry, buf *bytes.Buffer) error {

	var targetWidth, targetHeight = appConfig.Thumbmails.Width, appConfig.Thumbmails.Height

	decoded, format, err := image.Decode(bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		entry.WithError(err).Error("failed to decode image")
		return err
	}

	entry = entry.WithField("format", format)
	entry.Info("image decoded, generating thumbnail")

	dc := gg.NewContext(targetWidth, targetHeight)
	dc.SetRGB(0, 0, 0) // Set black background

	// Calculate the scaling factors to fit the image within the target dimensions
	scaleX := float64(appConfig.Thumbmails.Width) / float64(decoded.Bounds().Dx())
	scaleY := float64(appConfig.Thumbmails.Height) / float64(decoded.Bounds().Dy())
	scale := math.Min(scaleX, scaleY)

	// Calculate the dimensions for the scaled image
	scaledWidth := int(float64(decoded.Bounds().Dx()) * scale)
	scaledHeight := int(float64(decoded.Bounds().Dy()) * scale)

	// Calculate the position to center the scaled image on the canvas
	x := (targetWidth - scaledWidth) / 2
	y := (targetHeight - scaledHeight) / 2

	// Draw the scaled image onto the canvas with the calculated position
	dc.DrawImage(resize.Resize(uint(scaledWidth), uint(scaledHeight), decoded, resize.Lanczos3), x, y)

	var thumbnail = new(bytes.Buffer)
	err = dc.EncodePNG(thumbnail)
	if err != nil {
		entry.WithError(err).Error("failed to encode thumbnail")
		return err
	}

	entry.Info("thumbnail generated")
	entry.Info("uploading thumbnail")

	_, err = a.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(appConfig.Buckets.Photos),
		Key:         aws.String(fmt.Sprintf("thumbnails/%s", detail.Object.Key)),
		Body:        thumbnail,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		entry.WithError(err).Error("failed to upload thumbnail")
		return err
	}

	entry.Info("thumbnail uploaded")
	entry.Info("uploading original")

	_, err = a.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(appConfig.Buckets.Photos),
		Key:         aws.String(fmt.Sprintf("original/%s", detail.Object.Key)),
		Body:        bytes.NewBuffer(buf.Bytes()),
		ContentType: object.ContentType,
	})
	if err != nil {
		entry.WithError(err).Error("failed to upload thumbnail")
		return err
	}

	entry.Info("original uploaded")
	entry.Info("deleting staged image")

	_, err = a.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(appConfig.Buckets.Uploads),
		Key:    aws.String(detail.Object.Key),
	})
	if err != nil {
		entry.WithError(err).Error("failed to delete staged image")
		return err
	}

	entry.Info("staged image deleted")

	return nil

}

func main() {

	ctx := context.TODO()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to load aws default config")
	}

	loadConfig(ctx, awsCfg)

	dbx, _, err := initDB()
	if err != nil {
		logger.WithError(err).Error("failed to initialize database")
	}

	s3Client := s3.NewFromConfig(awsCfg)

	imageRepo := mysql.NewImageRepository(dbx)

	app := app{
		logger: logger,
		s3:     s3Client,
		image:  imageRepo,
	}

	lambda.Start(app.handle)

}
