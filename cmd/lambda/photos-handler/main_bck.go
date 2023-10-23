package main

// import (
// 	"fmt"
// 	"image"
// 	"log"
// 	"math"
// 	"os"

// 	_ "image/jpeg" // Import the image/jpeg package for JPEG support
// 	_ "image/png"  // Import the image/png package for PNG support

// 	"github.com/fogleman/gg"
// 	"github.com/nfnt/resize"
// )

// func createThumbnail(inputPath, outputPath string, targetWidth, targetHeight int) error {
// 	// Open the input image file
// 	file, err := os.Open(inputPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Decode the input image
// 	inputImage, _, err := image.Decode(file)
// 	if err != nil {
// 		return err
// 	}

// 	// Create a new image context for drawing the thumbnail with black bars
// 	dc := gg.NewContext(targetWidth, targetHeight)
// 	dc.SetRGB(0, 0, 0) // Set black background

// 	// Calculate the scaling factors to fit the image within the target dimensions
// 	scaleX := float64(targetWidth) / float64(inputImage.Bounds().Dx())
// 	scaleY := float64(targetHeight) / float64(inputImage.Bounds().Dy())
// 	scale := math.Min(scaleX, scaleY)

// 	// Calculate the dimensions for the scaled image
// 	scaledWidth := int(float64(inputImage.Bounds().Dx()) * scale)
// 	scaledHeight := int(float64(inputImage.Bounds().Dy()) * scale)

// 	// Calculate the position to center the scaled image on the canvas
// 	x := (targetWidth - scaledWidth) / 2
// 	y := (targetHeight - scaledHeight) / 2

// 	// Draw the scaled image onto the canvas with the calculated position
// 	dc.DrawImage(resize.Resize(uint(scaledWidth), uint(scaledHeight), inputImage, resize.Lanczos3), x, y)

// 	// Save the thumbnail to the output path
// 	if err := dc.SavePNG(outputPath); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func main() {
// 	inputPath := "input.jpg"      // Input image path
// 	outputPath := "thumbnail.jpg" // Output thumbnail path
// 	targetWidth := 512            // Desired width for the thumbnail
// 	targetHeight := 512           // Desired height for the thumbnail

// 	err := createThumbnail(inputPath, outputPath, targetWidth, targetHeight)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Thumbnail created successfully.")
// }
