package backupService

import (
	"bytes"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type UploadFileOptions struct {
	FolderId string
	Filepath string
}

func UploadFile(options UploadFileOptions) (*drive.File, error) {
	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("☁️ Uploading file to Google Drive...")
	}
	service, err := initDriveService()
	if err != nil {
		return nil, fmt.Errorf("unable to init drive service: %v", err)
	}

	folderID, err := getFolderID(service, os.Getenv("SITE_NAME"), options.FolderId)
	if err != nil {
		return nil, err
	}
	if folderID == "" {
		folderID, err = createFolder(service, os.Getenv("SITE_NAME"), options.FolderId)
		if err != nil {
			return nil, err
		}
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("📁 Parent Folder ID:", options.FolderId)
		fmt.Println("📁 Site Folder ID:", folderID)
		fmt.Println("📄 Filepath:", options.Filepath)
	}

	// Open the file
	localFile, err := os.Open(options.Filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer localFile.Close()

	// Detect the content type of the file
	contentType := mime.TypeByExtension(filepath.Ext(options.Filepath))
	if contentType == "" {
		// Default to plain text if type could not be detected
		contentType = "text/plain"
	}

	// Get the filename from the Filepath
	_, filename := filepath.Split(options.Filepath)

	// Create a file on Google Drive
	driveFile := &drive.File{
		Name:    filename,
		Parents: []string{folderID},
	}
	uploadedFile, err := service.Files.Create(driveFile).Media(localFile, googleapi.ContentType(contentType)).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Printf("✅ File '%s' uploaded with ID: %s\n", filename, uploadedFile.Id)
	}
	return uploadedFile, nil
}

type UploadBufferOptions struct {
	FolderId string
	Filename string
	Buffer   *bytes.Buffer
}

func UploadBuffer(options UploadBufferOptions) (*drive.File, error) {
	service, err := initDriveService()
	if err != nil {
		return nil, fmt.Errorf("unable to init drive service: %v", err)
	}

	folderID, err := getFolderID(service, os.Getenv("SITE_NAME"), options.FolderId)
	if err != nil {
		return nil, err
	}
	if folderID == "" {
		folderID, err = createFolder(service, os.Getenv("SITE_NAME"), options.FolderId)
		if err != nil {
			return nil, err
		}
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("📁 Parent Folder ID:", options.FolderId)
		fmt.Println("📁 Site Folder ID:", folderID)
		fmt.Println("📄 Filename:", options.Filename)
	}

	// Detect the content type of the file
	contentType := mime.TypeByExtension(filepath.Ext(options.Filename))
	if contentType == "" {
		// Default to plain text if type could not be detected
		contentType = "text/plain"
	}

	// Create a file on Google Drive
	driveFile := &drive.File{
		Name:    options.Filename,
		Parents: []string{folderID},
	}
	file, err := service.Files.Create(driveFile).Media(bytes.NewReader(options.Buffer.Bytes()), googleapi.ContentType(contentType)).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Printf("Buffer content '%s' uploaded with ID: %s\n", options.Filename, file.Id)
	}
	return file, nil
}
