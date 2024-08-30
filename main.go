package main

import (
	imagehandling "color-schemorator/image-handling"
	parsepalette "color-schemorator/parse-palette"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

const version = "1.0.0"

func main() {
	if slices.Contains(os.Args, "--version") {
		printVersionMessage()
		os.Exit(0)
	} else if len(os.Args) != 4 {
		printHelpMessage()
		os.Exit(1)
	}

	start := time.Now()

	colorSchemePath := os.Args[1]
	imgInputPath, err := validateExtension(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	imgOutputPath, err := validateExtension(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	palette, err := parsepalette.ParsePalette(colorSchemePath)
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

	err = imagehandling.SaveNewImg(imgOutputPath, newImg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	duration := time.Since(start)
	fmt.Printf("Completion Time: %v\n", duration)
}

func printVersionMessage() {
	fmt.Printf("Color Schemorator version %v\n", version)
	fmt.Println("Color Schemorator is a tool that adjusts the color palette of",
		"an image based on a provided list of hex color codes.")
	fmt.Println("For more information, visit: https://github.com/vannrr/color-schemorator")
}

func printHelpMessage() {
	fmt.Println("Error: Invalid number of arguments.")
	fmt.Println("Usage: csor <palette-file> <image-input> <image-output> [--version]")
	fmt.Println()
	fmt.Println("Description:")
	fmt.Println("  Color Schemorator is a tool that adjusts the color palette of an image")
	fmt.Println("  based on a provided list of hex color codes.")
	fmt.Println("  It creates a new image from the original colors")
	fmt.Println("  of the old image with their closest matches from the color")
	fmt.Println("  palette defined in the file (file can have '//' comments).")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <palette-file> Path to the plain text file containing hex color codes, one per line.")
	fmt.Println("  <image-input>  Path to the input image file (supported formats: jpg, jpeg, png).")
	fmt.Println("  <image-output> Path to the output image file (supported formats: jpg, jpeg, png).")
	fmt.Println("  --version      Display the version of the Color Schemorator tool.")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  csor colors.txt original-image.jpg new-image.jpg")
	fmt.Println()
	fmt.Println("Note: The order of arguments is important.")
}

func validateExtension(filePathString string) (string, error) {
	if ext := filepath.Ext(filePathString); ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("Invalid file extension for image file: '%v'", ext)
	} else {
		return filePathString, nil
	}
}
