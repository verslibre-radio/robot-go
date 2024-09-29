package main

import (
	"bytes"
	"fmt"
	"github.com/mjoes/mixcloud-go/pkg/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func FullShowName(metadata Metadata, date string) string {
  var name string
  if metadata.dj_name != "" {
      name = fmt.Sprintf("%s with %s (%s)", metadata.show_name, metadata.dj_name, date)
  } else {
      name = fmt.Sprintf("%s (%s)", metadata.show_name, date)
  }
  return name
}

func MixcloudUpload(srcPath string, localPicPath string, metadata Metadata, date string) error {
	audioFile, _ := os.Open(srcPath)
	picFile, _ := os.Open(localPicPath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	mp3Part, _ := writer.CreateFormFile("mp3", filepath.Base(srcPath))
	_, _ = io.Copy(mp3Part, audioFile)

	picPart, _ := writer.CreateFormFile("picture", filepath.Base(localPicPath))
	_, _ = io.Copy(picPart, picFile)
	writer.WriteField("name", FullShowName(metadata, date))
	writer.WriteField("description", metadata.description)
	writer.WriteField("hide_stats", "true")
	writer.WriteField("publish_date", utils.GetPublish())
	writer.WriteField("tags-0-tag", metadata.tags0)
	writer.WriteField("tags-1-tag", metadata.tags1)
	writer.WriteField("tags-2-tag", metadata.tags2)
	writer.WriteField("tags-3-tag", metadata.tags3)
	writer.WriteField("tags-4-tag", metadata.tags4)
	_ = writer.Close()

	url := fmt.Sprintf("https://api.mixcloud.com/upload/?access_token=%s", os.Getenv("API_KEY"))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Upload to Mixcloud %s failed: %v\n", metadata.show_name, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Upload to Mixcloud %s failed: %s\n", metadata.show_name, resp.Status)
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", responseBody)
		if bytes.Contains(responseBody, []byte("RateLimitException")) {
			return fmt.Errorf("RateLimit Exception, break program")
		}
	} else {
		log.Printf("Upload to Mixcloud %s PASSED\n", metadata.show_name)
	}
	resp.Body.Close()
	audioFile.Close()
	picFile.Close()

	return nil
}
