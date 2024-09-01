package utility

import (
	"fmt"
	"os"
	"path/filepath"
)

var validExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
}

func ValidateExtension(filePathString string, name string) error {
	ext := filepath.Ext(filePathString)
	if _, valid := validExtensions[ext]; !valid {
		return fmt.Errorf("invalid file extension on %v file: %v", name, ext)
	}

	return nil
}

func ValidateFileSize(file *os.File, name string, maxMB int64) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info: %w", err)
	}

	if fileInfo.Size() > maxMB*1024*1024 {
		return fmt.Errorf("%v file size %d bytes exceeds %dMB", name, fileInfo.Size(), maxMB)
	}

	return nil
}
