package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	backupService "github.com/CalebBarnes/wp-auto-backup/services/backup"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")

	fmt.Print("\033[32m") // Set color to green
	fmt.Print(`
 ________________________
< Starting WP Auto Backup >
 ------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
`)
	fmt.Print("\033[0m") // Reset color
	fmt.Println("")

	minutesStr := os.Getenv("BACKUP_INTERVAL_MINUTES")
	if minutesStr == "" {
		minutesStr = "1440"
	}
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		fmt.Println("Error converting BACKUP_INTERVAL_MINUTES to an integer:", err)
		return
	}

	// fmt.Println("Creating and uploading database dumps and wordpress directory backups every " + minutesStr + " minutes")
	fmt.Println("Backups enabled:")
	fmt.Println("- WP CLI Database dump: " + os.Getenv("SITE_NAME"))
	fmt.Println("- remote site directory: " + os.Getenv("REMOTE_SITE_DIR"))
	// fmt.Println("Frequency: " + minutesStr + " minutes")
	if minutes > 60 {
		hours := minutes / 60
		fmt.Println("- Frequency: " + strconv.Itoa(hours) + " hours")
	} else {
		fmt.Println("- Frequency: " + minutesStr + " minutes")
	}
	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("- Verbose: true")
	}
	fmt.Println("- Connecting to \033[4m" + os.Getenv("SSH_USER") + "@" + os.Getenv("SSH_HOST") + "\033[0m")
	fmt.Println("")

	if os.Getenv("GOOGLE_DRIVE_FOLDER_ID") != "" {
		backupService.UploadReadme(os.Getenv("GOOGLE_DRIVE_FOLDER_ID"))
	}

	if os.Getenv("BACKUP_ON_START") == "true" {
		runJob()
	}

	ticker := time.NewTicker(time.Minute * time.Duration(minutes))
	defer ticker.Stop()

	// Setting up a channel to listen for interrupt signal (Ctrl + C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Using a channel to communicate when to stop the loop
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ticker.C:

				fmt.Println("\n🚀Starting scheduled backup job at " + time.Now().Format("2006-01-02 15:04:05") + "\n")
				runJob()
			case <-sigs:
				fmt.Println("\nReceived an interrupt, stopping...")
				done <- true
				return
			}
		}
	}()
	// Wait for signal to stop
	<-done
	fmt.Println("Program exiting")
}

func runJob() {
	if os.Getenv("SSH_PORT") == "" {
		os.Setenv("SSH_PORT", "22")
	}

	println("")
	fmt.Println("🧙 Starting scheduled backup job at " + time.Now().Format("2006-01-02 15:04:05"))
	println("")

	user := os.Getenv("SSH_USER")
	host := os.Getenv("SSH_HOST")

	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02-150405")

	if os.Getenv("DATABASE_BACKUPS_DISABLED") != "true" {
		backupService.BackupDatabase(backupService.BackupDatabaseOptions{
			User: user,
			Host: host,
			Port: os.Getenv("SSH_PORT"),
		}, timestamp)
	}

	if os.Getenv("FILE_BACKUPS_DISABLED") != "true" {
		backupService.BackupFiles(backupService.BackupFilesOptions{
			User:                   user,
			Host:                   host,
			DownloadDestinationDir: "temp_files",
			ZipDestinationDir:      "backups",
		}, timestamp)
	}

	println("")
	fmt.Println("🧙‍♂️ Finished scheduled backup job at " + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Total time: " + time.Since(currentTime).String() + "🏃‍♂️💨⚡️")
	println("")
}
