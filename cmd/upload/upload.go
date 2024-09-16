package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjoes/mixcloud-go/pkg/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

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

	fmt.Println("Starting upload of audio")

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile(*cred_path))
	if err != nil {
		log.Fatal("Error:", err)
	}
	files, _ := os.ReadDir(audio_base_path)

	log.Println("Looping through files in the to be uploaded folder")
	for _, f := range files {
		tag := strings.Split(f.Name(), "_")[1]
		metadata := get_metadata(*db_path, tag)

		log.Println(f.Name(), "- Starting mixcloud upload process")
		log.Println(f.Name(), "- Downloading picture to local storage")
		picture_path := filepath.Join(picture_base_path, metadata.picture)
		utils.GetPicture(metadata.picture, driveService, picture_path, drive_picture_folder)

		log.Println(f.Name(), "- Start upload to mixcloud")
		audio_path := filepath.Join(audio_base_path, f.Name())
		err = MixcloudUpload(audio_path, picture_path, metadata)
		if err != nil {
			log.Fatal("Error:", err)
			return
		}
		log.Println(f.Name(), "- Start upload to Radiocult")
		log.Println("Yet to be implemented")

		log.Println(f.Name(), "- Start upload to Drive Archive")
		err = utils.Upload(driveService, f.Name(), audio_path, archive_id)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println(f.Name(), "- Move to local archive")
		err = utils.LocalMove(audio_path, filepath.Join(*archive_path, f.Name()))
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println(f.Name(), "- Update show nr")
		update_show_nr(*db_path, tag)
		log.Println("-----COMPLETED-----")
	}
	log.Println("Finished, check log for errors. Exiting program....")
}
