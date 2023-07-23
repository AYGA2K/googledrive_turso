package googledrive

import (
	"fmt"
	"io"
	"log"

	"google.golang.org/api/drive/v3"
)

func CreateDir(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}
	r, err := service.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	for _, i := range r.Files {
		fmt.Printf("%s (%s)\n", i.Name, i.Id)
		if i.Name == name {
			fmt.Println("directory already exists ")
			return i, nil
		}
	}

	file, err := service.Files.Create(d).Do()
	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}
	return file, nil
}

func CreateFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()
	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}
