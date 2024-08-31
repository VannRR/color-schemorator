package imagehandling

import (
	parsepalette "color-schemorator/parse-palette"
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

func ExtractPalette(inputImage image.Image) color.Palette {
	var colors sync.Map
	bounds := inputImage.Bounds()
	numCPU := runtime.NumCPU()

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
				r, g, b, a := inputImage.At(x, y).RGBA()
				rB := byte(r >> 8)
				gB := byte(g >> 8)
				bB := byte(b >> 8)
				aB := byte(a >> 8)
				rgba := color.RGBA{rB, gB, bB, aB}

				actual, _ := colors.LoadOrStore(rgba, uint32(1))
				if actual != uint32(1) {
					colors.Store(rgba, actual.(uint32)+1)
				}
			}
		}
	}

	for i := 0; i < strips; i++ {
		wg.Add(1)
		go processStrip(i)
	}

	wg.Wait()

	type kv struct {
		rgba  color.RGBA
		count uint32
	}

	var kvSlice []kv
	colors.Range(func(key, value interface{}) bool {
		kvSlice = append(kvSlice, kv{key.(color.RGBA), value.(uint32)})
		return true
	})

	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].count > kvSlice[j].count
	})

	var palette color.Palette
	for i := 0; i < parsepalette.MaxColors && i < len(kvSlice); i++ {
		palette = append(palette, kvSlice[i].rgba)
	}

	return palette
}

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
			endX = bounds.Max.X
		}
		wg.Add(1)
		go processStrip(startX, endX)
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
