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

func GetPaths(input_paths []string) (string, string, string, error) {
	if len(input_paths) < 3 {
		return "", "", "", fmt.Errorf("Too few arguments provided")
	} else if len(input_paths) > 3 {
		return "", "", "", fmt.Errorf("Too many arguments provided")
	}
	base_path := input_paths[1]
	audio_base_path := filepath.Join(base_path, "to_upload")
	picture_base_path := filepath.Join(base_path, "picture")
	CheckPath(audio_base_path)
	CheckPath(picture_base_path)

	cred_path := input_paths[2]
	CheckPath(cred_path)

	return audio_base_path, picture_base_path, cred_path, nil
}
