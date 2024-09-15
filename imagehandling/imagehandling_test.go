package imagehandling

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/VannRR/color-schemorator/parsepalette"
)

const testPaletteInputPath = "../parsepalette/test-palette-input.txt"
const testImgInputPath = "test-img-input.png"
const testImgOutputExpectPath = "test-img-output-expect.png"
const testImgOutputSavePath = "test-img-output-save" // add  ext in function

func Test_GetDecodedImage(t *testing.T) {
	_, err := GetDecodedImage(testImgInputPath)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}
}

func Test_ExtractPalette(t *testing.T) {
	width, height := 100, 100

	expectedPalette := color.Palette{
		color.RGBA{0xaf, 0xff, 0xff, 0xff},
		color.RGBA{0xbf, 0xff, 0xff, 0xff},
		color.RGBA{0xcf, 0xff, 0xff, 0xff},
		color.RGBA{0xdf, 0xff, 0xff, 0xff},
		color.RGBA{0xef, 0xff, 0xff, 0xff},
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	segment := height / len(expectedPalette)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := y / segment
			if i >= len(expectedPalette) {
				i = len(expectedPalette) - 1
			}
			img.Set(x, y, expectedPalette[i])
		}
	}

	actualPalette := ExtractPalette(img)

	if len(expectedPalette) != len(actualPalette) {
		t.Fatalf("Expected same palette len, got: expected=%v, actual=%v",
			len(expectedPalette), len(actualPalette))
	}

	for i := 0; i < len(expectedPalette); i++ {
		if !slices.Contains(expectedPalette, actualPalette[i]) {
			t.Errorf("expected palette: %v does not contain actualPalette[%v]=%v",
				expectedPalette, i, actualPalette[i])
		}
	}
}

func Test_CreateNewImg(t *testing.T) {
	oldImg, err := GetDecodedImage(testImgInputPath)
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}

	palette, err := parsepalette.ParsePalette(testPaletteInputPath)
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
