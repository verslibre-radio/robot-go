package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjoes/mixcloud-go/pkg/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var sheet_id string = "1XmJ8mXzMsBzDv13ZwM9tXasym5z3ZzlmNKC7xFudkzo"
var drive_picture_folder string = "1t7JgNd4U1oQEYw4NTdHPUFAIxd9YJWq3"

func main() {
	arg_list := os.Args
	audio_base_path, picture_base_path, cred_path, err := utils.GetPaths(arg_list)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	fmt.Println("Starting upload of audio")

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile(cred_path))
	if err != nil {
		log.Fatal("Error:", err)
	}
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsFile(cred_path))
	if err != nil {
		log.Fatal("Error:", err)
	}

	my_sheet, err := sheetsService.Spreadsheets.Values.Get(sheet_id, "meta").Do()
	if err != nil {
		log.Fatal("Error:", err)
	}

	files, _ := os.ReadDir(audio_base_path)
	columns := my_sheet.Values[0]

	log.Println("Looping through files in the to be uploaded folder")
	for _, f := range files {
		tag := strings.Split(f.Name(), "_")[1]

		log.Println(f.Name(), "- Starting mixcloud upload process")
		for _, sheet := range my_sheet.Values {
			if sheet[0] == tag {
				payload := make(map[string]string)
				for i := 0; i < len(sheet); i++ {
					payload[columns[i].(string)] = sheet[i].(string)
				}
				log.Println(f.Name(), "- Downloading picture to local storage")
				picture_path := filepath.Join(picture_base_path, sheet[4].(string))
				utils.GetPicture(sheet[4].(string), driveService, picture_path, drive_picture_folder)
				// audio_path := filepath.Join(audio_base_path, f.Name())

				log.Println(f.Name(), "- Start upload to mixcloud")
				// err = MixcloudUpload(audio_path, picture_path, payload)
				// if err != nil {
				// 	log.Fatal("Error:", err)
				// 	return
				// }
				log.Println("****", f.Name(), "- Start upload to Radiocult ***")
				log.Println(f.Name(), "- Start upload to Drive Archive")
				log.Println(f.Name(), "- Move to local archive")
				log.Println("-----COMPLETED-----")
			}
		}
	}
	log.Println("Finished, check log for errors. Exiting program....")
}
