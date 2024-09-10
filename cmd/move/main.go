package main

import (
	"context"
	"log"
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
	log.Println("Starting Google Drive move operation")
	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithCredentialsFile("./cred.json"))
	// copy_file("1fCDnR8KDTqwQEDEFuPZCvZqn8AInUiY_", "macmini", driveService)
	// move_file(driveService, "1iY1F3HvZulMLcLKjgnBoa4ZMKte9vl0Z", "sent", "macmini")
  utils.DownloadFile(driveService, "1pB9900KeU83_pT6iJ7uXwgTggf3XsR1b", "/Users/mortenslingsby/Desktop/test.mp3")
	// for _, f := range list_files(folderid_sent, driveService) {
	// 	fmt.Println(f.Id, f.Name)
	// }

}
