package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func RadiocultUpload(srcPath string, metadata Metadata) error {
	audioFile, _ := os.Open(srcPath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	mp3Part, err := writer.CreateFormFile("stationMedia", filepath.Base(srcPath))
  if err != nil {
    log.Printf("Failed to create form file part: %v\n", err)
    return err
  }
	_, _ = io.Copy(mp3Part, audioFile)
	_ = writer.Close()

	url := fmt.Sprintf("https://api.radiocult.fm/api/station/%s/media/track", os.Getenv("STATION_ID"))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-api-key", os.Getenv("RADIOCULT_API"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

  fmt.Println(resp.Status)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 201 {
		fmt.Printf("Upload to Radiocult %s failed: %s\n", metadata.show_name, resp.Status)
    responseBody, _ := io.ReadAll(resp.Body)
		if bytes.Contains(responseBody, []byte("RateLimitException")) {
			return fmt.Errorf("RateLimit Exception, break program")
		}
    return fmt.Errorf("Response: %s\n", responseBody)
	} 
  log.Printf("Upload to Radiocult %s PASSED\n", metadata.show_name)

	resp.Body.Close()
	audioFile.Close()

	return nil
}
