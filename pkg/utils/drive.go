package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

func ListFiles(folderid string, srv *drive.Service) []*drive.File {
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

func CopyFile(source_file_id string, dest_folder_id string, srv *drive.Service) {
	file, _ := srv.Files.Get(source_file_id).Do()
	copiedFile := &drive.File{
		Name:    file.Name,
		Parents: []string{dest_folder_id},
	}
	_, err := srv.Files.Copy(source_file_id, copiedFile).Do()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Copied", file.Name, "to folder: ", dest_folder_id)
	}
}

func MoveFile(srv *drive.Service, source_file_id string, source_folder_id string, dest_folder_id string) {
	_, err := srv.Files.Update(source_file_id, nil).
		RemoveParents(source_file_id).
		AddParents(dest_folder_id).
		Do()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("File %s moved from folder %s to folder %s", source_file_id, source_folder_id, dest_folder_id)
}

func DownloadFile(srv *drive.Service, source_file_id string, local_path string) {
	res, err := srv.Files.Get(source_file_id).Download()
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	outFile, err := os.Create(local_path)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File %s downloaded to %s", source_file_id, local_path)
}
