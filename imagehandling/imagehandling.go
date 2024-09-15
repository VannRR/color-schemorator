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
	"sort"
	"sync"

	"github.com/VannRR/color-schemorator/parsepalette"
	"github.com/VannRR/color-schemorator/utility"
)

const maxImgFileSizeMB = 15

// GetDecodedImage opens, validates, and decodes an image from the given file path.
// It returns the decoded image or an error if any step fails.
func GetDecodedImage(filePathString string) (image.Image, error) {
	if err := utility.ValidateExtension(filePathString, "input image"); err != nil {
		return nil, err
	}

	file, err := os.Open(filePathString)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePathString, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", filePathString, cerr)
		}
	}()

	if err := utility.ValidateFileSize(file, "Input image", maxImgFileSizeMB); err != nil {
		return nil, fmt.Errorf("file size validation failed for %s: %w", filePathString, err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image %s: %w", filePathString, err)
	}

	return img, nil
}

// ExtractPalette extracts the most common colors from an image, returning them as a color.Palette.
func ExtractPalette(inputImage image.Image) color.Palette {
	bounds := inputImage.Bounds()
	numCPU := runtime.NumCPU()
	stripWidth := (bounds.Max.X - bounds.Min.X) / numCPU

	type colorCount struct {
		color color.RGBA
		count uint32
	}

	processStrip := func(startX, endX int, colorMap map[color.RGBA]uint32, wg *sync.WaitGroup) {
		defer wg.Done()

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := startX; x < endX; x++ {
				rgba := rgbaAt(inputImage, x, y)
				colorMap[rgba]++
			}
		}
	}

	var wg sync.WaitGroup
	colorMaps := make([]map[color.RGBA]uint32, numCPU)

	for i := range colorMaps {
		colorMaps[i] = make(map[color.RGBA]uint32)
	}

	for i := 0; i < numCPU; i++ {
		startX := bounds.Min.X + i*stripWidth
		endX := startX + stripWidth
		if i == numCPU-1 {
			endX = bounds.Max.X // Handle the last strip to ensure full width is covered
		}

		wg.Add(1)
		go processStrip(startX, endX, colorMaps[i], &wg)
	}

	wg.Wait()

	finalColorMap := make(map[color.RGBA]uint32)
	for _, cmap := range colorMaps {
		for c, count := range cmap {
			finalColorMap[c] += count
		}
	}

	var colorCountSlice []colorCount
	for c, count := range finalColorMap {
		colorCountSlice = append(colorCountSlice, colorCount{color: c, count: count})
	}

	sort.Slice(colorCountSlice, func(i, j int) bool {
		return colorCountSlice[i].count > colorCountSlice[j].count
	})

	var palette color.Palette
	for i := 0; i < parsepalette.MaxColors && i < len(colorCountSlice); i++ {
		palette = append(palette, colorCountSlice[i].color)
	}

	return palette
}

// rgbaAt extracts the RGBA color at a given pixel location
func rgbaAt(img image.Image, x, y int) color.RGBA {
	r, g, b, a := img.At(x, y).RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// GenerateNewImg creates a new image by applying a color palette to the old image.
// The function leverages parallel processing to enhance performance.
func GenerateNewImg(oldImg image.Image, palette color.Palette) *image.Paletted {
	bounds := oldImg.Bounds()
	newImg := image.NewPaletted(bounds, palette)
	numCPU := runtime.NumCPU()

	stripWidth := (bounds.Max.X - bounds.Min.X) / numCPU
	var wg sync.WaitGroup

	processStrip := func(startX, endX int) {
		defer wg.Done()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := startX; x < endX; x++ {
				newImg.Set(x, y, palette.Convert(oldImg.At(x, y)))
			}
		}
	}

	for i := 0; i < numCPU; i++ {
		startX := bounds.Min.X + i*stripWidth
		endX := startX + stripWidth
		if i == numCPU-1 {
			endX = bounds.Max.X // Ensure the last strip covers the full width
		}

		wg.Add(1)
		go processStrip(startX, endX)
	}

	wg.Wait()
	return newImg
}

// SaveNewImg saves the new paletted image to the specified file path.
// It supports saving in JPEG or PNG formats.
func SaveNewImg(filePathString string, newImg *image.Paletted) error {
	outputFile, err := os.Create(filePathString)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePathString, err)
	}
	defer func() {
		if cerr := outputFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", filePathString, cerr)
		}
	}()

	ext := filepath.Ext(filePathString)
	switch ext {
	case ".jpg", ".jpeg":
		if err := jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 90}); err != nil {
			return fmt.Errorf("failed to encode image as JPEG: %w", err)
		}
	case ".png":
		if err := png.Encode(outputFile, newImg); err != nil {
			return fmt.Errorf("failed to encode image as PNG: %w", err)
		}
	default:
		return fmt.Errorf("invalid file extension for image output file: %v", ext)
	}

	return nil
}
