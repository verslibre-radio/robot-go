package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mjoes/mixcloud-go/pkg/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var folderIds = map[string]string{
	"auphonic_macmini":    "1ZbwEJnbv6OXJ3PF4cwMurElgHpZHClK7",
	"auphonic":            "1jX5SgOub7DKdyPUznNEGmEf0krf4fVyH",
	"auphonic_preprocess": "12Wn1XiyCDTAn1xI14CXuDTvT4aYl6DWr",
	"upload":              "1wGLWtfs4qEhHH_wtD2FPnEhp-n0s3haY",
	"upload_source":       "1My8d_fthYsRg0yV59kkkTtX6pktxFWu8",
	"sent":                "1tLBobAXugrZ5cHTP6X4Zk0LJDYnoZsN5",
	"macmini":             "1--7daLxjUi6zmFm6K952EDby8LBZFi5i",
}

var local_download_path = "/Users/mortenslingsby/Desktop/"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No path to local temp storage for files provided")
		return
	} else if len(os.Args) > 2 {
		fmt.Println("Too many arguments provided")
		return
	}

	base_path := os.Args[1]
	utils.CheckPath(base_path)

	log.Println("Starting Google Drive move operation")
	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithCredentialsFile("./cred.json"))

	// Auphonic on Macmini to auphonic upload
	log.Println("Processing the to be mixed auphonic files")
	for _, f := range utils.ListFiles(folderIds["auphonic_macmini"], driveService) {
		log.Println(f.Name, "- Copying...")
		utils.CopyFile(f.Id, folderIds["auphonic_preprocess"], driveService)
		log.Println(f.Name, "- Moving...")
		utils.MoveFile(driveService, f.Id, folderIds["auphonic_macmini"], folderIds["sent"])
		log.Println(f.Name, "- Successfully processed")
	}

	// Download from auphonic upload
	log.Println("Processing the mixed auphonic files")
	for _, f := range utils.ListFiles(folderIds["auphonic"], driveService) {
		log.Println(f.Name, "- Downloading...")
		utils.DownloadFile(driveService, f.Id, filepath.Join(base_path, f.Name))
		log.Println(f.Name, "- Moving...")
		utils.MoveFile(driveService, f.Id, folderIds["auphonic"], folderIds["sent"])
		log.Println(f.Name, "- Successfully processed")
	}

	// Download from upload
	log.Println("Processing the upload folder")
	for _, f := range utils.ListFiles(folderIds["upload_source"], driveService) {
		log.Println(f.Name, "- Downloading...")
		utils.DownloadFile(driveService, f.Id, filepath.Join(base_path, f.Name))
		log.Println(f.Name, "- Moving...")
		utils.MoveFile(driveService, f.Id, folderIds["upload_source"], folderIds["sent"])
		log.Println(f.Name, "- Successfully processed")
	}

	log.Println("Completed all the moves")
}