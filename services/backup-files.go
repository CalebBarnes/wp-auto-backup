package backupService

import (
	"fmt"
	"os"

	utils "github.com/Jambaree/wpe-backup-cloner/utils"
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
