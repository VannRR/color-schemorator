package parsepalette

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

const MinColors int = 2
const MaxColors int = 128
const maxParseErrors = 15
const maxPaletteFileSizeMB = 1

func ParsePalette(paletteInputPath string) (color.Palette, error) {
	file, err := os.Open(paletteInputPath)
	if err != nil {
		return color.Palette{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return color.Palette{}, err
	}

	if fileInfo.Size() > maxPaletteFileSizeMB*1024*1024 {
		return color.Palette{},
			fmt.Errorf("Palette file size exceeds %vMB\n", maxPaletteFileSizeMB)
	}

	byteSlice, err := io.ReadAll(file)
	if err != nil {
		return color.Palette{}, err
	}

	lines := strings.Split(string(byteSlice), "\n")
	var linesFiltered []string
	for i := 0; i < len(lines); i++ {
		line := strings.Split(lines[i], "//")[0]
		line = strings.TrimSpace(line)
		if line != "" {
			linesFiltered = append(linesFiltered, line)
		}
	}
	lines = linesFiltered

	colors, err := parseColorsFromLines(lines)
	if err != nil {
		return color.Palette{}, err
	}

	return colors, nil
}

func parseColorsFromLines(lines []string) ([]color.Color, error) {
	colors := make([]color.Color, 0, MaxColors)
	errors := make([]string, 0, maxParseErrors)
	errCount := 0
	for ln, line := range lines {
		rgba, err := parseHexColor(line)
		if len(colors) > MaxColors {
			errors = append(
				[]string{fmt.Sprintf("Max amount of colors in palette are %v", MaxColors)},
				errors...,
			)
			errCount += 1
			break
		}
		if err != nil {
			if errCount < maxParseErrors {
				errors = append(errors, fmt.Sprintf("Error on line %v: %v", ln+1, err))
			}
			errCount += 1
		}
		if !slices.Contains(colors, rgba) {
			colors = append(colors, rgba)
		}
	}

	if len(colors) < MinColors {
		errors = append(
			[]string{fmt.Sprintf("Minimum amount of colors in palette are %v", MinColors)},
			errors...,
		)
		errCount += 1
	}

	if errCount > 0 {
		allErrors := strings.Join(errors, "\n")
		if errCount >= maxParseErrors {
			allErrors = fmt.Sprintf(
				"%v\n%v more errors...", allErrors, errCount-len(errors))
		}
		return colors, fmt.Errorf(allErrors)
	}

	return colors, nil
}

func parseHexColor(hexColorString string) (color.Color, error) {
	if len(hexColorString) != 4 && len(hexColorString) != 7 {
		return color.RGBA{}, fmt.Errorf("Invalid hex color '%v'", truncateString(hexColorString, 30))
	}
	if !strings.HasPrefix(hexColorString, "#") {
		return color.RGBA{}, fmt.Errorf("Invalid hex color '%v'", truncateString(hexColorString, 30))
	}

	for _, c := range hexColorString[1:] {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return color.RGBA{}, fmt.Errorf("Invalid hex color '%v'", truncateString(hexColorString, 30))
		}
	}

	var r, g, b byte
	var err error
	if len(hexColorString) == 4 {
		r, err = parseHexPair(hexColorString[1:2] + hexColorString[1:2])
		g, err = parseHexPair(hexColorString[2:3] + hexColorString[2:3])
		b, err = parseHexPair(hexColorString[3:4] + hexColorString[3:4])
	} else if len(hexColorString) == 7 {
		r, err = parseHexPair(hexColorString[1:3])
		g, err = parseHexPair(hexColorString[3:5])
		b, err = parseHexPair(hexColorString[5:7])
	}
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

func truncateString(s string, length int) string {
	if len(s)+3 > length {
		s = s[0:length] + "..."
	}
	return s
}

func parseHexPair(hexPair string) (byte, error) {
	value, err := strconv.ParseInt(hexPair, 16, 9)
	if err != nil {
		return 0, err
	}
	return byte(value), nil
}

func SaveNewPalette(paletteOutputPath string, palette color.Palette) error {
	var hexColors []string
	for _, c := range palette {
		r, g, b, _ := c.RGBA()
		color := fmt.Sprintf("#%02X%02X%02X",
			byte(r>>8), byte(g>>8), byte(b>>8))
		hexColors = append(hexColors, color)
	}

	outputFile, err := os.Create(paletteOutputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	fileInfo, err := outputFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() > maxPaletteFileSizeMB*1024*1024 {
		return fmt.Errorf("Palette file size exceeds %vMB\n", maxPaletteFileSizeMB)
	}

	outputFile.WriteString(strings.Join(hexColors, "\n"))

	return nil
}
