package imagehandling

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const maxImgFileSizeMB = 15

func GetDecodedImage(filePathString string) (image.Image, error) {
	ext := filepath.Ext(filePathString)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, fmt.Errorf("Invalid file extension on input file: %v", ext)
	}
	file, err := os.Open(filePathString)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() > maxImgFileSizeMB*1024*1024 {
		return nil,
			fmt.Errorf("Input image file size exceeds %vMB\n", maxImgFileSizeMB)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil,
			fmt.Errorf("Error decoding image: %v", err)
	}

	return img, nil
}

func GenerateNewImg(oldImg image.Image, palette color.Palette) *image.Paletted {
	bounds := oldImg.Bounds()
	newImg := image.NewPaletted(bounds, palette)
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	strips := numCPU
	stripWidth := (bounds.Max.X - bounds.Min.X) / strips

	var wg sync.WaitGroup

	processStrip := func(stripIndex int) {
		defer wg.Done()

		startX := bounds.Min.X + stripIndex*stripWidth
		endX := startX + stripWidth
		if stripIndex == strips-1 {
			endX = bounds.Max.X
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := startX; x < endX; x++ {
				newImg.Set(x, y, palette.Convert(oldImg.At(x, y)))
			}
		}
	}

	for i := 0; i < strips; i++ {
		wg.Add(1)
		go processStrip(i)
	}

	wg.Wait()

	return newImg
}

func SaveNewImg(filePathString string, newImg *image.Paletted) error {
	outputFile, err := os.Create(filePathString)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer outputFile.Close()

	fileInfo, err := outputFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() > maxImgFileSizeMB*1024*1024 {
		return fmt.Errorf("Output image file size exceeds %vMB\n", maxImgFileSizeMB)
	}

	ext := filepath.Ext(filePathString)

	if ext == ".jpg" || ext == ".jpeg" {
		if err := jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("Error encoding new image: %v", err)
		}
	} else if ext == ".png" {
		if err := png.Encode(outputFile, newImg); err != nil {
			return fmt.Errorf("Error encoding new image: %v", err)
		}
	} else {
		return fmt.Errorf("Invalid file extension for image output file: %v", ext)
	}

	return nil
}
