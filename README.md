# Color Schemorator

```
Usage:
  csor -m generate -p <palettePath> -i <imgInputPath> -o <imgOutputPath>
  csor -m extract -i <imgInputPath> -P <paletteOutputPath>
  csor -v
  csor -h

Description:
  Color Schemorator modifies an image's color palette based on a given list
  of hex color codes (file can have '//' comments) or extracts the color
  palette from an image.

  - Generate mode: Creates a new image by replacing its colors with the
    closest matches from the specified palette.
  - Extract mode: Extracts the color palette from an image (in order of
    occurrence) and saves it to a file.

Arguments:
  -m   Mode of operation: 'generate' or 'extract'.
  -p   Path to the plain text file containing hex color codes, one per line
       (required for 'generate' mode).
  -i   Path to the input image file (supported formats: jpg, jpeg, png).
  -o   Path to the output image file (supported formats: jpg, jpeg, png)
       (required for 'generate' mode).
  -P   Path to the output palette file (required for 'extract' mode).
  -v   Display the version of the Color Schemorator tool.
  -h   Display this help message.

Example:
  csor -m generate -p colors.txt -i original-image.jpg -o new-image.jpg
  csor -m extract -i original-image.jpg -P palette.txt
```

## Install

To install the Color Schemorator application,
follow these steps (note `go` and `make` need to be installed first):

1) Clone the repository:
```
git clone https://github.com/vannrr/color-schemorator.git
cd color-schemorator
```

2) Build and install:
```
make install
```
By default, the application will be installed to ~/bin/.
If you wish to change the installation directory,
you can modify the INSTALL_DIR variable in the Makefile:
```
INSTALL_DIR := /your/custom/path/
```
3) Verify the installation:
```
 ~/bin/csor -v
```
Replace ~/bin/ with your custom path if you changed the INSTALL_DIR.

## Example Images

| Before | After |
|--------|-------|
| <img src="/image-handling/test-img-input.png?raw=true" width="300"/> | <img src="/image-handling/test-img-output-expect.png?raw=true" width="300"/> |

## License

This project is licensed under the MIT License. See the LICENSE file for details.
