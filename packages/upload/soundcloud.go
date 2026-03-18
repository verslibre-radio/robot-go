package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"robot-go/utils"
)

const soundcloudAPIBaseURL = "https://api.soundcloud.com"
const soundcloudUploadLimitBytes int64 = 500 * 1024 * 1024

type soundcloudTrackResponse struct {
	URN string `json:"urn"`
}

func SoundcloudUpload(srcPath string, localPicPath string, metadata Metadata, date string, accessToken string) (string, error) {
	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		return "", err
	}
	if fileInfo.Size() > soundcloudUploadLimitBytes {
		return "", fmt.Errorf("SoundCloud upload limit exceeded: %s is %d bytes", srcPath, fileInfo.Size())
	}

	trackURN, err := createSoundcloudTrack(srcPath, metadata, date, accessToken)
	if err != nil {
		return "", err
	}
	if err := updateSoundcloudTrack(trackURN, srcPath, localPicPath, metadata, date, accessToken); err != nil {
		return "", err
	}

	return trackURN, nil
}

func createSoundcloudTrack(srcPath string, metadata Metadata, date string, accessToken string) (string, error) {
	audioFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer audioFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	audioPart, err := writer.CreateFormFile("track[asset_data]", filepath.Base(srcPath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(audioPart, audioFile); err != nil {
		return "", err
	}

	if err := writer.WriteField("track[title]", FullShowName(metadata, date)); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tracks", soundcloudAPIBaseURL), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", accessToken))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SoundCloud track upload failed: %s: %s", resp.Status, string(responseBody))
	}

	var track soundcloudTrackResponse
	if err := json.Unmarshal(responseBody, &track); err != nil {
		return "", fmt.Errorf("failed to parse SoundCloud track response: %w", err)
	}
	if track.URN == "" {
		return "", fmt.Errorf("SoundCloud upload succeeded but no track URN was returned")
	}

	log.Printf("Upload to SoundCloud %s PASSED\n", metadata.show_name)
	return track.URN, nil
}

func updateSoundcloudTrack(trackURN string, srcPath string, localPicPath string, metadata Metadata, date string, accessToken string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("track[title]", FullShowName(metadata, date)); err != nil {
		return err
	}
	if err := writer.WriteField("track[description]", metadata.description); err != nil {
		return err
	}
	if err := writer.WriteField("track[tag_list]", soundcloudTagList(metadata)); err != nil {
		return err
	}
	if err := writer.WriteField("track[sharing]", "public"); err != nil {
		return err
	}
	if err := writer.WriteField("track[publish_at]", utils.GetPublish()); err != nil {
		return err
	}

	if localPicPath != "" {
		picFile, err := os.Open(localPicPath)
		if err != nil {
			return err
		}
		defer picFile.Close()

		picPart, err := writer.CreateFormFile("track[artwork_data]", filepath.Base(localPicPath))
		if err != nil {
			return err
		}
		if _, err := io.Copy(picPart, picFile); err != nil {
			return err
		}
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tracks/%s", soundcloudAPIBaseURL, trackURN), body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", accessToken))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SoundCloud metadata update failed for %s: %s: %s", filepath.Base(srcPath), resp.Status, string(responseBody))
	}

	return nil
}

func soundcloudTagList(metadata Metadata) string {
	tags := []string{metadata.tags0, metadata.tags1, metadata.tags2, metadata.tags3, metadata.tags4}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if strings.ContainsRune(tag, ' ') {
			result = append(result, fmt.Sprintf("\"%s\"", tag))
			continue
		}
		result = append(result, tag)
	}
	return strings.Join(result, " ")
}
