package parsepalette

import (
	"image/color"
	"testing"
)

const testPaletteInputPath = "test-palette-input.txt"
const testPaletteOutputPath = "test-palette-output.txt"

func Test_ParsePalette(t *testing.T) {
	expectedColors := color.Palette{
		color.RGBA{51, 51, 51, 255},
		color.RGBA{63, 54, 86, 255},
		color.RGBA{220, 138, 120, 255},
		color.RGBA{221, 120, 120, 255},
		color.RGBA{234, 118, 203, 255},
		color.RGBA{136, 57, 239, 255},
		color.RGBA{210, 15, 57, 255},
		color.RGBA{230, 69, 83, 255},
		color.RGBA{254, 100, 11, 255},
		color.RGBA{223, 142, 29, 255},
		color.RGBA{64, 160, 43, 255},
		color.RGBA{23, 146, 153, 255},
		color.RGBA{4, 165, 229, 255},
		color.RGBA{32, 159, 181, 255},
		color.RGBA{30, 102, 245, 255},
		color.RGBA{114, 135, 253, 255},
		color.RGBA{76, 79, 105, 255},
		color.RGBA{92, 95, 119, 255},
		color.RGBA{108, 111, 133, 255},
		color.RGBA{124, 127, 147, 255},
		color.RGBA{140, 143, 161, 255},
		color.RGBA{156, 160, 176, 255},
		color.RGBA{172, 176, 190, 255},
		color.RGBA{188, 192, 204, 255},
		color.RGBA{204, 208, 218, 255},
		color.RGBA{239, 241, 245, 255},
		color.RGBA{230, 233, 239, 255},
		color.RGBA{220, 224, 232, 255},
	}
	palette, err := ParsePalette(testPaletteInputPath)
	if err != nil {
		t.Errorf("Expected no error, got error: %v", err)
	}

	for i, e := range expectedColors {
		if e != palette[i] {
			t.Errorf("Expected color %v, got %v", e, palette[i])
		}
	}
}

func Test_ParseHexColor(t *testing.T) {
	tests := []struct {
		hexColorString string
		expected       color.RGBA
		isError        bool
	}{
		{"#fff", color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 255}, false},
		{"#ffffff", color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 255}, false},
		{"#000", color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}, false},
		{"#000000", color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 255}, false},
		{"#123abc", color.RGBA{R: 0x12, G: 0x3a, B: 0xbc, A: 255}, false},
		{"123456", color.RGBA{}, true},
		{"#12g456", color.RGBA{}, true},
		{"#1234", color.RGBA{}, true},
	}

	for _, tt := range tests {
		rgba, err := parseHexColor(tt.hexColorString)

		if err != nil && !tt.isError {
			t.Errorf("Expected no error for input %v, but got: %v", tt.hexColorString, err)
		} else if err == nil && tt.isError {
			t.Errorf("Expected error for input %v, but got none", tt.hexColorString)
		} else if rgba != tt.expected {
			t.Errorf("Expected color %v for input, but got %v", tt.expected, rgba)
		}
	}
}

func Test_ParseHexPair(t *testing.T) {
	tests := []struct {
		hexPair  string
		expected byte
		hasError bool
	}{
		{"00", 0x00, false},
		{"7F", 0x7F, false},
		{"FF", 0xFF, false},
		{"01", 0x01, false},
		{"0A", 0x0A, false},
		{"G1", 0x00, true},  // Invalid hex character
		{"", 0x00, true},    // Empty string
		{"123", 0x00, true}, // More than two characters
	}

	for _, tt := range tests {
		result, err := parseHexPair(tt.hexPair)
		if tt.hasError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got nil", tt.hexPair)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for input %q, but got %v", tt.hexPair, err)
			}
			if result != tt.expected {
				t.Errorf("Expected %x for input %q, but got %x", tt.expected, tt.hexPair, result)
			}
		}
	}
}

func Test_SaveNewPalette(t *testing.T) {
	writePalette := color.Palette{
		color.RGBA{0xaf, 0xff, 0xff, 0xff},
		color.RGBA{0xbf, 0xff, 0xff, 0xff},
		color.RGBA{0xcf, 0xff, 0xff, 0xff},
		color.RGBA{0xdf, 0xff, 0xff, 0xff},
		color.RGBA{0xef, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}

	err := SaveNewPalette(testPaletteOutputPath, writePalette)
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}

	readPalette, err := ParsePalette(testPaletteOutputPath)
	if err != nil {
		t.Fatalf("Expected no error, got error: %v", err)
	}

	if len(writePalette) != len(readPalette) {
		t.Fatalf("Expected same palette len, got: write=%v, read=%v",
			len(writePalette), len(readPalette))
	}

	for i := 0; i < len(writePalette); i++ {
		if writePalette[i] != readPalette[i] {
			t.Errorf("expected: %v, got: %v", writePalette[i], readPalette[i])
		}
	}
}
