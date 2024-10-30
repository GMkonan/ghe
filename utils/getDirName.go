package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetDirName() string {
	currentDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Error grabing the path to dir")
	}
	dirName := filepath.Base(currentDir)

	// make sure dirName is valid for github
	if strings.Contains(dirName, " ") {
		dirName = strings.ReplaceAll(dirName, " ", "-")
	}
	return dirName
}
