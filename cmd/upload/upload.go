package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mjoes/mixcloud-go/pkg/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var sheet_id string = "1XmJ8mXzMsBzDv13ZwM9tXasym5z3ZzlmNKC7xFudkzo"
var drive_picture_folder string = "1t7JgNd4U1oQEYw4NTdHPUFAIxd9YJWq3"
var archive_id string = "1qklZQWVpNRYJWLd0-0zBhxZLyWCHrtpe"

func main() {
	base_path := flag.String("local", "/var/lib/robot", "Path to local temp storage for upload files and pictures")
	archive_path := flag.String("archive", "", "Path to local archive folder")
	cred_path := flag.String("credentials", "/etc/robot/cred.json", "Path to credentials file")
	db_path := flag.String("metadata", "/var/lib/robot/metadata.sql", "Path to credentials file")

	flag.Parse()
	if *archive_path == "" {
		fmt.Println("Archive path not set")
		return
	}

	audio_base_path, picture_base_path, err := utils.GetPaths(*base_path)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Open connection to metadata DB
	sqlDB, _ := sql.Open("sqlite3", *db_path)
	defer sqlDB.Close()

	fmt.Println("Starting upload of audio")

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile(*cred_path))
	if err != nil {
		log.Fatal("Error:", err)
	}
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsFile(*cred_path))
	if err != nil {
		log.Fatal("Error:", err)
	}
	sheet_meta, err := sheetsService.Spreadsheets.Values.Get(sheet_id, "meta").Do()
	if err != nil {
		log.Fatal("Error:", err)
	}

	files, _ := os.ReadDir(audio_base_path)

	log.Println("Looping through files in the to be uploaded folder")
	for _, f := range files {
		if string(f.Name()[0]) == "." {
			continue
		}
		split_name := strings.Split(f.Name(), "_")
		date := split_name[0]
		tag := split_name[1]

		metadata := get_metadata(sheet_meta, sqlDB, tag, date)
		fmt.Println(metadata)
		new_meta_row(sqlDB, date, metadata)

		log.Println(f.Name(), "- Starting mixcloud upload process")
		log.Println(f.Name(), "- Downloading picture to local storage")
		picture_path := filepath.Join(picture_base_path, metadata.picture)
		utils.GetPicture(metadata.picture, driveService, picture_path, drive_picture_folder)
		audio_path := filepath.Join(audio_base_path, f.Name())

		// Mixcloud
		log.Println(f.Name(), "- Start upload to mixcloud")
		if get_meta_status(sqlDB, "mixcloud", tag, date) {
			err = MixcloudUpload(audio_path, picture_path, metadata)
			if err != nil {
				log.Fatal("Error:", err)
				return
			}
			update_meta_status(sqlDB, "mixcloud", tag, date)
		} else {
			log.Println("File already uploaded to Mixcloud")
		}

		// Radiocult
		log.Println(f.Name(), "- Start upload to Radiocult")
		if get_meta_status(sqlDB, "radiocult", tag, date) {
			err = RadiocultUpload(audio_path, metadata)
			if err != nil {
				log.Fatal("Error:", err)
				return
			}
      update_meta_status(sqlDB, "radiocult", tag, date)
		} else {
			log.Println("File already uploaded to Radiocult")
		}

		// Google drive archive
		log.Println(f.Name(), "- Start upload to Drive Archive")
		if get_meta_status(sqlDB, "drive", tag, date) {
			err = utils.Upload(driveService, f.Name(), audio_path, archive_id)
			if err != nil {
			  log.Fatal(err)
			  return
			}
			update_meta_status(sqlDB, "drive", tag, date)
		} else {
			log.Println("File already uploaded to Drive")
		}

		// Local drive archive
		log.Println(f.Name(), "- Move to local archive")
		if ready_for_upload(sqlDB, tag, date) {
			err = utils.LocalMove(audio_path, filepath.Join(*archive_path, f.Name()))
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println(f.Name(), "- Update show nr")
			update_show_nr(sqlDB, tag)
		} else {
			log.Println("Not all upload stages complete, not moving to archive")
		}
		log.Println("-----COMPLETED-----")
	}
	log.Println("Finished, check log for errors. Exiting program....")
}
