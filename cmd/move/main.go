package main

import (
	"fmt"
	// "google.golang.org/api/drive/v3"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)

var folderid_auphonic_macmini string = "1ZbwEJnbv6OXJ3PF4cwMurElgHpZHClK7"
var folderid_auphonic string = "1jX5SgOub7DKdyPUznNEGmEf0krf4fVyH"
var folderid_auphonic_preprocess string = "12Wn1XiyCDTAn1xI14CXuDTvT4aYl6DWr"
var folderid_upload string = "1wGLWtfs4qEhHH_wtD2FPnEhp-n0s3haY"
var folderid_upload_source string = "1My8d_fthYsRg0yV59kkkTtX6pktxFWu8"
var folderid_sent string = "1tLBobAXugrZ5cHTP6X4Zk0LJDYnoZsN5"

func list_files(folderid string) {
	driveService, _ := drive.NewService(ctx, option.WithCredentialsFile("./cred.json"))
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

func main() {
	fmt.Println("Hello")
}
