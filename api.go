package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func mixcloud_upload(srcPath string, localPicPath string, fileName string, payload map[string]string) {
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
	writer.WriteField("disable_comments", "true")
	writer.WriteField("hide_stats", "true")
	writer.WriteField("unlisted", "true")
	writer.WriteField("publish_date", get_publish())

	for key, value := range payload {
		if strings.Contains(key, "tags") {
			writer.WriteField(key, value)
		}
	}
	_ = writer.Close()

	url := fmt.Sprintf("https://api.mixcloud.com/upload/?access_token=%s", os.Getenv("API_KEY"))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Upload to Mixcloud %s failed: %v\n", fileName, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Upload to Mixcloud %s failed: %s\n", fileName, resp.Status)
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", responseBody)
		if bytes.Contains(responseBody, []byte("RateLimitException")) {
			fmt.Println("RateLimit Exception, break program")
			return
		}
	} else {
		fmt.Printf("Upload to Mixcloud %s PASSED\n", fileName)
	}
	resp.Body.Close()
	audioFile.Close()
}
