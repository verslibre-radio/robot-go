package main

import (
	"context"
	"fmt"
	"log"

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

func list_files(folderid string, srv *drive.Service) []*drive.File {
	var Q_string string = fmt.Sprintf("'%s' in parents", folderid)
	drive := srv.Files.List()
	filtered, err := drive.
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Corpora("drive").
		DriveId("0AGvEMGW0880aUk9PVA").
		Q(Q_string).
		Do()

	if err != nil {
		log.Println("Error:", err)
	}
	return filtered.Files
}

func copy_file(source_file_id string, dest_folder string, srv *drive.Service) {
	dest_folder_id := folderIds[dest_folder]
	file, _ := srv.Files.Get(source_file_id).Do()
	copiedFile := &drive.File{
		Name:    file.Name,
		Parents: []string{dest_folder_id},
	}
	_, err := srv.Files.Copy(source_file_id, copiedFile).Do()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Copied", file.Name, "to folder: ", dest_folder)
	}
}

func move_file(srv *drive.Service, source_file_id string, source_folder string, dest_folder string) {
	_, err := srv.Files.Update(source_file_id, nil).
		RemoveParents(folderIds[source_folder]).
		AddParents(folderIds[dest_folder]).
		Do()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("File %s moved from folder %s to folder %s", source_file_id, source_folder, dest_folder)
}

func main() {
	log.Println("Starting Google Drive move operation")
	ctx := context.Background()
	driveService, _ := drive.NewService(ctx, option.WithCredentialsFile("./cred.json"))
	// copy_file("1fCDnR8KDTqwQEDEFuPZCvZqn8AInUiY_", "macmini", driveService)
	move_file(driveService, "1iY1F3HvZulMLcLKjgnBoa4ZMKte9vl0Z", "sent", "macmini")
	// for _, f := range list_files(folderid_sent, driveService) {
	// 	fmt.Println(f.Id, f.Name)
	// }

}
