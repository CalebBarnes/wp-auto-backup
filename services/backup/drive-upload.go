package backupService

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type UploadFileOptions struct {
	FolderId string
	Filepath string
}

func UploadFileInSiteFolder(options UploadFileOptions) (*drive.File, error) {
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
		fmt.Println("📁 Site Folder ID:", folderID)
	}

	return UploadFile(UploadFileOptions{
		FolderId: folderID,
		Filepath: options.Filepath,
	})
}

func UploadFile(options UploadFileOptions) (*drive.File, error) {
	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("☁️ Uploading file to Google Drive...")
	}
	service, err := initDriveService()
	if err != nil {
		return nil, fmt.Errorf("unable to init drive service: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("📁 Parent Folder ID:", options.FolderId)
		fmt.Println("📄 Filepath:", options.Filepath)
	}

	// Open the file
	localFile, err := os.Open(options.Filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer localFile.Close()

	// Get the file size and create a progress reader
	fileInfo, err := localFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file info: %v", err)
	}
	progressReader, err := UploadProgressReader(localFile, fileInfo.Size(), func(readSize int64, totalSize int64, speed float64) {
		uploadedMB := float64(readSize) / (1024 * 1024)
		totalMB := float64(totalSize) / (1024 * 1024)
		percentage := float64(readSize) / float64(totalSize) * 100
		fmt.Printf("📤 Uploading: %.2fMB/%.2fMB (%.2f%%) at %.2fMB/s\r", uploadedMB, totalMB, percentage, speed)
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create progress reader: %v", err)
	}

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
	uploadedFile, err := service.Files.Create(driveFile).Media(progressReader, googleapi.ContentType(contentType)).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Printf("✅ File '%s' uploaded with ID: %s\n", filename, uploadedFile.Id)
	}
	return uploadedFile, nil
}

func UploadBufferInSiteFolder(options UploadBufferOptions) (*drive.File, error) {
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
		fmt.Println("📁 Site Folder ID:", folderID)
	}

	return UploadBuffer(UploadBufferOptions{
		FolderId: folderID,
		Filename: options.Filename,
		Buffer:   options.Buffer,
	})
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

	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("📁 Parent Folder ID:", options.FolderId)
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
		Parents: []string{options.FolderId},
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

type ReportFunc func(int64, int64, float64)
type ProgressReader struct {
	reader     io.Reader
	totalSize  int64
	readSize   int64
	startTime  time.Time
	reportFunc ReportFunc
}

func UploadProgressReader(reader io.Reader, totalSize int64, reportFunc ReportFunc) (io.Reader, error) {
	return &ProgressReader{
		reader:     reader,
		totalSize:  totalSize,
		startTime:  time.Now(),
		reportFunc: reportFunc,
	}, nil
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.readSize += int64(n)
	if pr.reportFunc != nil {
		elapsed := time.Since(pr.startTime).Seconds()
		speed := float64(pr.readSize) / elapsed / (1024 * 1024) // Speed in MB/s
		pr.reportFunc(pr.readSize, pr.totalSize, speed)
	}
	return n, err
}
