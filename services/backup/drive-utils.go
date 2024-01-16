package backupService

import (
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
)

func UploadReadme(folderId string) (bool, error) {
	service, err := initDriveService()
	if err != nil {
		return false, fmt.Errorf("failed to initialize drive service: %v", err)
	}
	query := fmt.Sprintf("name='readme.txt' and '%s' in parents", folderId)

	response, err := service.Files.List().Q(query).Do()
	if err != nil {
		return false, fmt.Errorf("failed to query file: %v", err)
	}

	for _, file := range response.Files {
		if file.Name == "readme.txt" {
			return true, nil // Return true if readme.txt exists
		}
	}

	// readme.txt does not exist, so create it
	createdFile, err := UploadFile(UploadFileOptions{
		FolderId: folderId,
		Filepath: "example/readme.txt",
	})

	if err != nil {
		return false, fmt.Errorf("failed to upload file: %v", err)
	}
	if os.Getenv("VERBOSE") == "true" {
		fmt.Printf("Created file %s (%s)\n", createdFile.Name, createdFile.Id)
	}
	return false, nil // readme.txt does not exist
}

func getFolderID(service *drive.Service, folderName string, parentFolderId string) (string, error) {
	var query string
	if parentFolderId != "" {
		query = fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and '%s' in parents", folderName, parentFolderId)
	} else {
		query = fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s'", folderName)
	}

	response, err := service.Files.List().Q(query).Do()
	if err != nil {
		return "", fmt.Errorf("failed to query folder: %v", err)
	}

	for _, file := range response.Files {
		return file.Id, nil // Return the first matching folder
	}

	return "", nil // No folder found
}

func createFolder(service *drive.Service, folderName string, parentFolderId string) (string, error) {
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}
	if parentFolderId != "" {
		folder.Parents = []string{parentFolderId}
	}
	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create folder: %v", err)
	}
	return createdFolder.Id, nil
}
