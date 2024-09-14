package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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
	}
}

func MoveFile(srv *drive.Service, source_file_id string, source_folder_id string, dest_folder_id string) {
	_, err := srv.Files.Update(source_file_id, nil).
		RemoveParents(source_folder_id).
		AddParents(dest_folder_id).
		Do()

	if err != nil {
		log.Fatal(err)
	}
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
}

func CheckPath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return os.MkdirAll(absPath, os.ModePerm)
	}
	return nil
}

func GetPicture(tag string, driveService *drive.Service, picture_path string, drive_picture_folder string) {
	var Q_string string = fmt.Sprintf("'%s' in parents", drive_picture_folder)
	drive := driveService.Files.List()
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
	for _, f := range filtered.Files {
		if f.Name == tag {
			DownloadFile(driveService, f.Id, picture_path)
		}
	}
}

func Upload(srv *drive.Service, file_name string, file_path string, destination_folder string) error {
	file := &drive.File{
		Name:     file_name,
		MimeType: "audio/mpeg",
		Parents:  []string{destination_folder},
	}

	f, err := os.Open(file_path)
	if err != nil {
		return fmt.Errorf("Unable to open file: %v", err)
	}
	defer f.Close()

	_, err = srv.Files.
		Create(file).
		SupportsAllDrives(true).
		Media(f).
		Do()
	if err != nil {
		return fmt.Errorf("Unable to create file: %v", err)
	}
	return nil
}
