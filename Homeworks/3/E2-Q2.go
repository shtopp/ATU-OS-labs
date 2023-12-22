package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"sync"
)

var (
	kernel = [][]float64{
		{1.0 / 16, 2.0 / 16, 1.0 / 16},
		{2.0 / 16, 4.0 / 16, 2.0 / 16},
		{1.0 / 16, 2.0 / 16, 1.0 / 16},
	}
)

func main() {
	// Open the original image file
	file, err := os.Open("ATU.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the original image
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	outputImage := image.NewRGBA(img.Bounds())

	// Create a channel to pass image chunks between goroutines
	chunkSize := 100 // Adjust the chunk size as needed
	chunkChan := make(chan image.Image, runtime.NumCPU())
	blurredChunkChan := make(chan image.Image, runtime.NumCPU())
	final := make(chan struct{})

	// Start goroutines to blurring image chunks
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for imgChunk := range chunkChan {
				blurredChunk := applyGaussianBlur(imgChunk)
				// Send blurred chunk to the next stage or process it further
				blurredChunkChan <- blurredChunk
				fmt.Print(".")
			}
			//wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		//close(blurredChunkChan)
	}()

	go func() {
		defer close(final)
		for blurredChunk := range blurredChunkChan {
			//copy the blurred chunks to the final image
			draw.Draw(outputImage, blurredChunk.Bounds(), blurredChunk, image.Point{}, draw.Src)
		}
	}()

	// Split the image into chunks and send them to the channel for processing
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += chunkSize {
		chunkBounds := image.Rect(bounds.Min.X, y, bounds.Max.X, y+chunkSize)
		imgChunk := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(chunkBounds)
		chunkChan <- imgChunk
	}
	close(chunkChan) // Close the channel to signal that all chunks have been sent

	// Wait for all goroutines to finish processing
	wg.Wait()
	close(blurredChunkChan)
	<-final

	// Save the final blurred image to file
	outputFile, err := os.Create("output_blurred.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, outputImage) // Replace 'img' with the final blurred image
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Blurred image generated.")
}

// Apply Gaussian blur to an image chunk
func applyGaussianBlur(imgChunk image.Image) image.Image {
	fmt.Print("G")
	bounds := imgChunk.Bounds()
	blurredImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Apply the kernel to each pixel
			var r, g, b, a float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky

					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						pr, pg, pb, pa := imgChunk.At(px, py).RGBA()
						weight := kernel[ky+1][kx+1]

						r += float64(pr) * weight
						g += float64(pg) * weight
						b += float64(pb) * weight
						a += float64(pa) * weight
					}
				}
			}

			// Set the blurred pixel values
			blurredImg.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
		}
	}

	return blurredImg
}
