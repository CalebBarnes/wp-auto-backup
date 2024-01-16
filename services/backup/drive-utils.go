package backupService

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

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
