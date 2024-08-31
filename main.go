package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var sheet_id string = "1XmJ8mXzMsBzDv13ZwM9tXasym5z3ZzlmNKC7xFudkzo"

func main() {
	fmt.Println("hello world")

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsFile("./cred.json"))
	if err != nil {
		log.Println("Error:", err)
	}
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("./cred.json"))

	if err != nil {
		log.Println("Error:", err)
	}

	my_sheet, err := sheetsService.Spreadsheets.Values.Get(sheet_id, "meta").Do()
	if err != nil {
		log.Println("Error:", err)
	}
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
		fmt.Println(f.Id, f.Name)
	}

	files, _ := os.ReadDir("./upload")

	for _, f := range files {
		tag := strings.Split(f.Name(), "_")[1]
		for _, sheet := range my_sheet.Values {
			if sheet[0] == tag {
				fmt.Println(tag)
			}
		}
	}
}
