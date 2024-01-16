package backupService

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var service *drive.Service

func InitDriveService() {
	// Read the JSON key file of your service account
	ctx := context.Background()
	serviceAccountKeyFile := os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_FILE")
	b, err := os.ReadFile(serviceAccountKeyFile)
	if err != nil {
		log.Fatalf("Unable to read service account key file: %v", err)
	}

	// Authenticate using the service account
	driveConfig, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse service account key file to config: %v", err)
	}
	client := driveConfig.Client(ctx)

	service, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
}

type UploadFileOptions struct {
	FolderId string
	Filepath string
}

func UploadFile(options UploadFileOptions) error {
	fmt.Println("☁️ Uploading file to Google Drive...")
	if service == nil {
		InitDriveService()
	}
	// Open the file
	localFile, err := os.Open(options.Filepath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
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
		Parents: []string{options.FolderId},
	}
	uploadedFile, err := service.Files.Create(driveFile).Media(localFile, googleapi.ContentType(contentType)).Do()
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}

	fmt.Printf("✅ File '%s' uploaded with ID: %s\n", filename, uploadedFile.Id)
	return nil
}

type UploadBufferOptions struct {
	FolderId string
	Filename string
	Buffer   *bytes.Buffer
}

func UploadBuffer(options UploadBufferOptions) error {
	if service == nil {
		InitDriveService()
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
		Parents: []string{options.FolderId},
	}
	file, err := service.Files.Create(driveFile).Media(bytes.NewReader(options.Buffer.Bytes()), googleapi.ContentType(contentType)).Do()
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}

	fmt.Printf("Buffer content '%s' uploaded with ID: %s\n", options.Filename, file.Id)
	return nil
}
