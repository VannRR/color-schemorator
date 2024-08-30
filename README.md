# Color Schemorator

```
Usage: csor <palette-file> <image-input> <image-output> [--version]

Description:
  Color Schemorator is a tool that adjusts the color palette of an image
  based on a provided list of hex color codes.
  It creates a new image from the original colors
  of the old image with their closest matches from the color
  palette defined in the file (file can have '//' comments).

Arguments:
  <palette-file> Path to the plain text file containing hex color codes, one per line.
  <image-input>  Path to the input image file (supported formats: jpg, jpeg, png).
  <image-output> Path to the output image file (supported formats: jpg, jpeg, png).
  --version      Display the version of the Color Schemorator tool.

Example:
  csor colors.txt original-image.jpg new-image.jpg

Note: The order of arguments is important.
```

## Install

To install the Color Schemorator application, follow these steps:

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
 ~/bin/csor --version
```
Replace ~/bin/ with your custom path if you changed the INSTALL_DIR.

## Example Images

| Before | After |
|--------|-------|
| <img src="/image-handling/test-img-input.png?raw=true" width="300"/> | <img src="/image-handling/test-img-output-expect.png?raw=true" width="300"/> |

## License

This project is licensed under the MIT License. See the LICENSE file for details.
