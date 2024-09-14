package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func GetPublish() string {
	now := time.Now()
	futureTime := now.Add(24 * time.Hour)
	format_date := fmt.Sprintf("%d-%02d-%02dT14:00:00Z", futureTime.Year(), futureTime.Month(), futureTime.Day())
	return format_date
}

func GetPaths(base_path string) (string, string, error) {
	audio_base_path := filepath.Join(base_path, "to_upload")
	picture_base_path := filepath.Join(base_path, "picture")
	CheckPath(audio_base_path)
	CheckPath(picture_base_path)

	return audio_base_path, picture_base_path, nil
}

func LocalMove(src_path string, dest_path string) error {
	srcFile, err := os.Open(src_path)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dest_path)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	err = os.Remove(src_path)
	if err != nil {
		return fmt.Errorf("failed to delete source file: %w", err)
	}

	return nil
}
