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
	"strings"
)

func MixcloudUpload(srcPath string, localPicPath string, payload map[string]string) error {
	audioFile, _ := os.Open(srcPath)
	picFile, _ := os.Open(localPicPath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	mp3Part, _ := writer.CreateFormFile("mp3", filepath.Base(srcPath))
	_, _ = io.Copy(mp3Part, audioFile)

	picPart, _ := writer.CreateFormFile("picture", filepath.Base(localPicPath))
	_, _ = io.Copy(picPart, picFile)
	writer.WriteField("name", payload["show_name"])
	writer.WriteField("description", payload["description"])
	writer.WriteField("hide_stats", "true")
	writer.WriteField("publish_date", utils.GetPublish())

	for key, value := range payload {
		if strings.Contains(key, "tags") {
			writer.WriteField(key, value)
		}
	}
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
		log.Printf("Upload to Mixcloud %s failed: %v\n", payload["show_name"], err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Upload to Mixcloud %s failed: %s\n", payload["show_name"], resp.Status)
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", responseBody)
		if bytes.Contains(responseBody, []byte("RateLimitException")) {
			return fmt.Errorf("RateLimit Exception, break program")
		}
	} else {
		log.Printf("Upload to Mixcloud %s PASSED\n", payload["show_name"])
	}
	resp.Body.Close()
	audioFile.Close()
	picFile.Close()

	return nil
}
