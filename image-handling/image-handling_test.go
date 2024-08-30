package imagehandling

import (
	"color-schemorator/parse-palette"
	"image"
	"image/color"
	"testing"
)

const testColorSchemePath = "../parse-palette/test-palette.txt"
const testImgInputPath = "test-img-input.png"
const testImgOutputExpectPath = "test-img-output-expect.png"
const testImgOutputSavePath = "test-img-output-save" // add  ext in function

func Test_GetDecodedImage(t *testing.T) {
	_, err := GetDecodedImage(testImgInputPath)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}
}

func Test_CreateNewImg(t *testing.T) {
	oldImg, err := GetDecodedImage(testImgInputPath)
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}

	palette, err := parsepalette.ParsePalette(testColorSchemePath)
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}

	newImgResult := GenerateNewImg(oldImg, palette)

	newImgExpect, err := GetDecodedImage(testImgOutputExpectPath)
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}

	if !compareImages(newImgResult, newImgExpect) {
		t.Errorf("Generated image does not match the expected image.")
	}
}

func compareImages(img1 *image.Paletted, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1 != bounds2 {
		return false
	}

	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			if img1.At(x, y) != img2.At(x, y) {
				return false
			}
		}
	}

	return true
}

func Test_SaveNewImg(t *testing.T) {
	width, height := 100, 100

	img := image.NewPaletted(image.Rect(0, 0, width, height), color.Palette{
		color.RGBA{255, 255, 255, 255},
	})

	fillColor := color.RGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, fillColor)
		}
	}

	err := SaveNewImg(testImgOutputSavePath+".jpg", img)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}
	err = SaveNewImg(testImgOutputSavePath+".jpeg", img)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}
	err = SaveNewImg(testImgOutputSavePath+".png", img)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}
}
