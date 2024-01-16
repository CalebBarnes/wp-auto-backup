package backupService

import (
	"fmt"
	"log"
	"os"

	utils "github.com/CalebBarnes/wp-auto-backup/utils"
)

type BackupFilesOptions struct {
	User                   string
	Host                   string
	DownloadDestinationDir string
	ZipDestinationDir      string
}

func BackupFiles(options BackupFilesOptions, timestamp string) {
	if options.ZipDestinationDir == "" {
		fmt.Println("Zip destination directory is required")
		return
	}
	// check if zip destination dir exists, if not create it
	if _, err := os.Stat(options.ZipDestinationDir); os.IsNotExist(err) {
		fmt.Println("Creating destination directory: " + options.ZipDestinationDir)
		err := os.MkdirAll(options.ZipDestinationDir, 0755)
		if err != nil {
			log.Fatalf("Error creating destination directory: %v", err)
			return
		}
	}

	fmt.Println("🗂️ Starting file backup...")
	err := utils.RsyncFromServer(utils.RsyncOptions{
		User:           options.User,
		Host:           options.Host,
		DestinationDir: options.DownloadDestinationDir,
		Verbose:        os.Getenv("VERBOSE") == "true",
	})
	if err != nil {
		fmt.Println("Error in rsync while backing up files:", err)
		return
	}

	zipFileName := fmt.Sprintf("%s/%s-wordpress-files-backup-%s.zip", options.ZipDestinationDir, os.Getenv("SITE_NAME"), timestamp)
	utils.CreateZipFile(zipFileName, options.DownloadDestinationDir)

	UploadFile(UploadFileOptions{
		FolderId: os.Getenv("GOOGLE_DRIVE_FOLDER_ID"),
		Filepath: zipFileName,
	})
}
