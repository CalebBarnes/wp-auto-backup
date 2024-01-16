package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	wpeService "github.com/Jambaree/wpe-backup-cloner/services"
)

func main() {
	godotenv.Load(".env.local")

	client := &http.Client{}
	runScheduledJob(client)

	// ticker := time.NewTicker(5 * time.Second)
	// defer ticker.Stop()

	// Setting up a channel to listen for interrupt signal (Ctrl + C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Using a channel to communicate when to stop the loop
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			// case <-ticker.C:
			// 	runScheduledJob(client)
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

func runScheduledJob(client *http.Client) {
	println("Running clone backup job...")

	sites, err := wpeService.GetSites(client)
	if err != nil {
		fmt.Println("Error getting sites:", err)
		return
	}

	// log the site name
	for _, site := range sites.Results {
		for _, install := range site.Installs {
			installNames := strings.Split(os.Getenv("WPE_INSTALLS"), ",")

			for _, name := range installNames {
				if install.Name == strings.TrimSpace(name) {
					fmt.Println("Found install to backup:", install.Name)
					fmt.Println("Found site to backup:", site.Name)
					handleBackup(client, install)
					// If you want to stop checking after finding a match, uncomment the following line
					// break
				}
			}
		}
	}
}

func handleBackup(client *http.Client, install wpeService.Install) {
	fmt.Println("Creating backup for install:", install.Name)

	// create the backup
	// backup, err := wpeService.CreateBackup(client, install)
	// if err != nil {
	// 	fmt.Println("Error creating backup:", err)
	// 	return
	// }
	// fmt.Println("Backup created:", backup.ID)
	// fmt.Println("Backup status:", backup.Status)

	// timeout := time.After(5 * time.Minute) // 5 minute timeout
	// tick := time.Tick(3 * time.Second)     // 3 second interval between backup status checks
	// for {
	// 	select {
	// 	case <-timeout:
	// 		fmt.Println("Backup check timed out.")
	// 		return
	// 	case <-tick:
	// 		backup, err := wpeService.GetBackupStatus(client, install, "673367b6-0451-48b2-ae86-0f7592e1bdf9")
	// 		if err != nil {
	// 			fmt.Println("Error getting backup status:", err)
	// 			return
	// 		}

	// 		fmt.Println("Backup status:", backup.Status)

	// 		if backup.Status == "completed" {
	// 			// Call your new function here
	// 			handleCloneBackup(backup)
	// 			return
	// 		}
	// 	}
	// }
}

// func handleCloneBackup(backup wpeService.BackupResponse) {
// 	fmt.Println("Cloning backup...")
// }
