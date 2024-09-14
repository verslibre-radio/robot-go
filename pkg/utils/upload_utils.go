package utils

import (
	"fmt"
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
