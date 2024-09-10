package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var sheet_id string = "1XmJ8mXzMsBzDv13ZwM9tXasym5z3ZzlmNKC7xFudkzo"

func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err) // or log.Print, log.Fatal, or handle in another way
	}
}

func get_picture(tag string, ctx context.Context) {
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("./cred.json"))
	drive := driveService.Files.List()
	filtered, err := drive.
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Corpora("drive").
		DriveId("0AGvEMGW0880aUk9PVA").
		Q("'1t7JgNd4U1oQEYw4NTdHPUFAIxd9YJWq3' in parents").
		Do()

	if err != nil {
		log.Println("Error:", err)
	}
	for _, f := range filtered.Files {
		if f.Name == tag {
			fmt.Println(f.Id, f.Name)
		}
	}
}

func get_publish() string {
	now := time.Now()
	futureTime := now.Add(24 * time.Hour)
	format_date := fmt.Sprintf("%d-%02d-%02dT14:00:00Z", futureTime.Year(), futureTime.Month(), futureTime.Day())
	return format_date
}

func main() {
	fmt.Println("Starting upload of audio")

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsFile("./cred.json"))
	if err != nil {
		log.Println("Error:", err)
	}

	my_sheet, err := sheetsService.Spreadsheets.Values.Get(sheet_id, "meta").Do()
	if err != nil {
		log.Println("Error:", err)
	}

	files, _ := os.ReadDir("./upload")

	columns := my_sheet.Values[0]
	for _, f := range files {
		tag := strings.Split(f.Name(), "_")[1]
		for _, sheet := range my_sheet.Values {
			if sheet[0] == tag {
				payload := make(map[string]string)
				for i := 0; i < len(sheet); i++ {
					payload[columns[i].(string)] = sheet[i].(string)
				}
				get_picture(sheet[4].(string), ctx)
				audio_path := filepath.Join("upload", f.Name())
				pic_path := filepath.Join("pictures", sheet[4].(string))
				mixcloud_upload(audio_path, pic_path, payload)
			}
		}
	}
}
