package main

import (
	"context"
	"encoding/json"
	"fmt"
	"photos/internal/store/mysql"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"

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
	image  *mysql.ImageRepository
}

func (a *app) handle(ctx context.Context, event *events.CloudWatchEvent) error {

	if event.Version != "0" {
		return fmt.Errorf("unsupported event version: %s", event.Version)
	}

	if event.Source != "aws.s3" {
		return fmt.Errorf("unsupported event source: %s", event.Source)
	}

	var detail = new(EventDetail)

	err := json.Unmarshal(event.Detail, detail)
	if err != nil {
		return err
	}

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

	imageRepo := mysql.NewImageRepository(dbx)

	app := app{
		logger: logger,
		image:  imageRepo,
	}

	lambda.Start(app.handle)

}
