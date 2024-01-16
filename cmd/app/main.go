package main

import (
	"fmt"
	"time"

	backupService "github.com/Jambaree/wpe-backup-cloner/services"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting up...")
	godotenv.Load(".env.local")

	runJob()

	// ticker := time.NewTicker(24 * time.Hour)
	// defer ticker.Stop()

	// // Setting up a channel to listen for interrupt signal (Ctrl + C)
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// // Using a channel to communicate when to stop the loop
	// done := make(chan bool, 1)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			runJob()
	// 		case <-sigs:
	// 			fmt.Println("\nReceived an interrupt, stopping...")
	// 			done <- true
	// 			return
	// 		}
	// 	}
	// }()
	// // Wait for signal to stop
	// <-done
	// fmt.Println("Program exiting")
}

func runJob() {
	user := "realcedar"
	host := "realcedar2dev.ssh.wpengine.net"
	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02-150405")
	fmt.Println(timestamp)

	backupService.BackupFiles(backupService.BackupFilesOptions{
		User:                   user,
		Host:                   host,
		DownloadDestinationDir: "temp_files/" + user,
		ZipDestinationDir:      "backups",
	}, timestamp)

	// backupService.BackupDatabase(backupService.BackupDatabaseOptions{
	// 	User: user,
	// 	Host: host,
	// 	Port: "22",
	// }, timestamp)
}
