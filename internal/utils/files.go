package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	merrors "github.com/tekofx/crossposter/internal/errors"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadFile(downloadURL, saveAs string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(saveAs)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}
	return nil
}

func ReadFile(name string) ([]byte, *merrors.MError) {
	fileBytes, err := os.ReadFile(name)
	if err != nil {
		merrors.New(merrors.CannotReadFileErrorCode, err.Error())
	}

	return fileBytes, nil
}
