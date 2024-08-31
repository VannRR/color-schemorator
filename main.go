package main

import (
	imagehandling "color-schemorator/image-handling"
	parsepalette "color-schemorator/parse-palette"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const version = "1.1.0"

func main() {
	versionFlag := flag.Bool("v", false, "Display the version of the Color Schemorator tool")
	helpFlag := flag.Bool("h", false, "Display help message")
	mode := flag.String("m", "", "Mode of operation: 'generate' or 'extract'")
	paletteInput := flag.String("p", "",
		"Path to the plain text file containing hex color codes, one per line (required for 'generate' mode)")
	imageInput := flag.String("i", "",
		"Path to the input image file (supported formats: jpg, jpeg, png)")
	imageOutput := flag.String("o", "",
		"Path to the output image file (supported formats: jpg, jpeg, png) (required for 'generate' mode)")
	paletteOutput := flag.String("P", "", "Path to the output palette file (required for 'extract' mode)")

	flag.Parse()

	if *versionFlag {
		printVersionMessage()
		os.Exit(0)
	} else if *helpFlag {
		printHelpMessage()
		os.Exit(0)
	}

	switch *mode {
	case "generate":
		if *paletteInput == "" || *imageInput == "" || *imageOutput == "" {
			printHelpMessage()
			os.Exit(1)
		}
		start := time.Now()
		generate(*paletteInput, *imageInput, *imageOutput)
		fmt.Println("Image generated successfully in", time.Since(start))

	case "extract":
		if *imageInput == "" || *paletteOutput == "" {
			printHelpMessage()
			os.Exit(1)
		}
		start := time.Now()
		extract(*imageInput, *paletteOutput)
		fmt.Println("Palette extracted successfully in", time.Since(start))

	default:
		printHelpMessage()
		os.Exit(1)
	}
}

func generate(paletteInputPath, imgInputPath, imgOutputPath string) {
	if err := validateExtension(imgInputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := validateExtension(imgOutputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	palette, err := parsepalette.ParsePalette(paletteInputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	oldImg, err := imagehandling.GetDecodedImage(imgInputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	newImg := imagehandling.GenerateNewImg(oldImg, palette)

	if err = imagehandling.SaveNewImg(imgOutputPath, newImg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func extract(imgInputPath, paletteOutputPath string) {
	if err := validateExtension(imgInputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	inputImg, err := imagehandling.GetDecodedImage(imgInputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	palette := imagehandling.ExtractPalette(inputImg)

	if err := parsepalette.SaveNewPalette(paletteOutputPath, palette); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func validateExtension(filePathString string) error {
	if ext := filepath.Ext(filePathString); ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("Invalid file extension for image file: '%v'", ext)
	} else {
		return nil
	}
}

func printVersionMessage() {
	fmt.Printf("Color Schemorator version %v\n", version)
	fmt.Println("Color Schemorator is a tool that adjusts the color palette of an image based")
	fmt.Println("on a provided list of hex color codes.")
	fmt.Println("For more information, visit:")
	fmt.Println("  https://github.com/vannrr/color-schemorator")
}

func printHelpMessage() {
	fmt.Println("Usage:")
	fmt.Println("  csor -m generate -p <palettePath> -i <imgInputPath> -o <imgOutputPath>")
	fmt.Println("  csor -m extract -i <imgInputPath> -P <paletteOutputPath>")
	fmt.Println("  csor -v")
	fmt.Println("  csor -h")
	fmt.Println()
	fmt.Println("Description:")
	fmt.Println("  Color Schemorator modifies an image's color palette based on a given list")
	fmt.Println("  of hex color codes (file can have '//' comments) or extracts the color")
	fmt.Println("  palette from an image.")
	fmt.Println()
	fmt.Println("  - Generate mode: Creates a new image by replacing its colors with the")
	fmt.Println("    closest matches from the specified palette.")
	fmt.Println("  - Extract mode: Extracts the color palette from an image (in order of")
	fmt.Println("    occurrence) and saves it to a file.")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  -m   Mode of operation: 'generate' or 'extract'.")
	fmt.Println("  -p   Path to the plain text file containing hex color codes, one per line")
	fmt.Println("       (required for 'generate' mode).")
	fmt.Println("  -i   Path to the input image file (supported formats: jpg, jpeg, png).")
	fmt.Println("  -o   Path to the output image file (supported formats: jpg, jpeg, png)")
	fmt.Println("       (required for 'generate' mode).")
	fmt.Println("  -P   Path to the output palette file (required for 'extract' mode).")
	fmt.Println("  -v   Display the version of the Color Schemorator tool.")
	fmt.Println("  -h   Display this help message.")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  csor -m generate -p colors.txt -i original-image.jpg -o new-image.jpg")
	fmt.Println("  csor -m extract -i original-image.jpg -P palette.txt")
}
