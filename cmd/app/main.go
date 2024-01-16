package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	backupService "github.com/CalebBarnes/wp-auto-backup/services"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("WP Auto Backup")
	godotenv.Load(".env.local")

	minutesStr := os.Getenv("BACKUP_INTERVAL_MINUTES")
	if minutesStr == "" {
		minutesStr = "1440"
	}
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		fmt.Println("Error converting BACKUP_INTERVAL_MINUTES to an integer:", err)
		return
	}

	fmt.Println("Running backup every " + minutesStr + " minutes")

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
	user := os.Getenv("SSH_USER")
	host := os.Getenv("SSH_HOST")

	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02-150405")

	backupService.BackupDatabase(backupService.BackupDatabaseOptions{
		User: user,
		Host: host,
		Port: "22",
	}, timestamp)

	backupService.BackupFiles(backupService.BackupFilesOptions{
		User:                   user,
		Host:                   host,
		DownloadDestinationDir: "temp_files/" + os.Getenv("SITE_NAME"),
		ZipDestinationDir:      "backups",
	}, timestamp)
}
