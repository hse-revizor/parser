package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalFileContent struct {
	Content string `json:"content"`
}

func ConvertLocalPathToFilePath(localPath string) (string, error) {
	if !filepath.IsAbs(localPath) {
		absPath, err := filepath.Abs(localPath)
		if err != nil {
			return "", fmt.Errorf("failed to convert relative path to absolute: %w", err)
		}
		localPath = absPath
	}
	if info, err := os.Stat(localPath); err != nil || info.IsDir() {
		return "", errors.New("path must point to a valid file")
	}

	return localPath, nil
}

func FetchFileContentFromLocalPath(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return string(content), nil
}
