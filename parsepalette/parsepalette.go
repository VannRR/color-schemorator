package parsepalette

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/VannRR/color-schemorator/utility"
)

const (
	MinColors            int = 2
	MaxColors            int = 128
	maxParseErrors           = 15
	maxPaletteFileSizeMB     = 1
)

// ParsePalette reads a palette file from the given path, validates its size,
// and parses the colors, returning a color.Palette.
func ParsePalette(paletteInputPath string) (color.Palette, error) {
	file, err := os.Open(paletteInputPath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	if err := utility.ValidateFileSize(file, "Input palette", maxPaletteFileSizeMB); err != nil {
		return nil, err
	}

	lines, err := readNonEmptyLines(file)
	if err != nil {
		return nil, err
	}

	colors, err := parseColorsFromLines(lines)
	if err != nil {
		return nil, err
	}

	return colors, nil
}

// readNonEmptyLines reads all non-empty lines from a file,
// ignoring comments, and returns them as a slice of strings.
func readNonEmptyLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(strings.Split(scanner.Text(), "//")[0])
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

// parseColorsFromLines parses the hex colors from the lines to RGBA,
// returning a slice of colors
func parseColorsFromLines(lines []string) ([]color.Color, error) {
	colors := make([]color.Color, 0, MaxColors)
	errors := make([]string, 0, maxParseErrors)
	errCount := 0
	seenColors := make(map[color.Color]struct{})

	for ln, line := range lines {
		rgba, err := parseHexColor(line)
		if err != nil {
			if errCount < maxParseErrors {
				errors = append(errors, fmt.Sprintf("Error on line %v: %v", ln+1, err))
			}
			errCount++
			if errCount >= maxParseErrors {
				break
			}
			continue
		}
		if len(colors) >= MaxColors {
			errors = append([]string{fmt.Sprintf("Max amount of colors in palette is %v", MaxColors)}, errors...)
			break
		}
		if _, exists := seenColors[rgba]; !exists {
			colors = append(colors, rgba)
			seenColors[rgba] = struct{}{}
		}
	}

	if len(colors) < MinColors {
		errors = append([]string{fmt.Sprintf("Minimum amount of colors in palette is %v", MinColors)}, errors...)
		errCount++
	}

	if errCount > 0 {
		allErrors := strings.Join(errors, "\n")
		if errCount > len(errors) {
			allErrors = fmt.Sprintf("%v\n%v more errors...", allErrors, errCount-len(errors))
		}
		return colors, fmt.Errorf(allErrors)
	}

	return colors, nil
}

// parseHexColor parses the RGBA color from a hex color string then returns it
func parseHexColor(hexColorString string) (color.Color, error) {
	hexColorString = strings.TrimSpace(hexColorString)
	if len(hexColorString) != 4 && len(hexColorString) != 7 {
		return color.RGBA{}, fmt.Errorf("invalid hex color '%v'", truncateString(hexColorString, 30))
	}
	if !strings.HasPrefix(hexColorString, "#") {
		return color.RGBA{}, fmt.Errorf("invalid hex color '%v'", truncateString(hexColorString, 30))
	}

	var r, g, b byte
	var err error
	if len(hexColorString) == 4 { // #RGB format
		r, err = parseHexPair(hexColorString[1:2] + hexColorString[1:2])
		if err == nil {
			g, err = parseHexPair(hexColorString[2:3] + hexColorString[2:3])
		}
		if err == nil {
			b, err = parseHexPair(hexColorString[3:4] + hexColorString[3:4])
		}
	} else if len(hexColorString) == 7 { // #RRGGBB format
		r, err = parseHexPair(hexColorString[1:3])
		if err == nil {
			g, err = parseHexPair(hexColorString[3:5])
		}
		if err == nil {
			b, err = parseHexPair(hexColorString[5:7])
		}
	}
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}

func truncateString(s string, length int) string {
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}

func parseHexPair(hexPair string) (byte, error) {
	value, err := strconv.ParseInt(hexPair, 16, 9)
	if err != nil {
		return 0, fmt.Errorf("invalid hex pair '%s': %w", hexPair, err)
	}
	return byte(value), nil
}

// SaveNewPalette saves a Palette as a plain text file of hex colors,
// one color per line
func SaveNewPalette(paletteOutputPath string, palette color.Palette) error {
	hexColors := make([]string, 0, len(palette))
	for _, c := range palette {
		r, g, b, _ := c.RGBA()
		color := fmt.Sprintf("#%02X%02X%02X", byte(r>>8), byte(g>>8), byte(b>>8))
		hexColors = append(hexColors, color)
	}

	outputFile, err := os.Create(paletteOutputPath)
	if err != nil {
		return fmt.Errorf("failed to create palette file: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	_, err = writer.WriteString(strings.Join(hexColors, "\n"))
	if err != nil {
		return fmt.Errorf("failed to write palette to file: %w", err)
	}

	return nil
}
