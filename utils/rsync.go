package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

type RsyncOptions struct {
	User           string
	Host           string
	DestinationDir string
	Verbose        bool
}

const (
	initialDelay  = 5 * time.Second
	maxRetries    = 5
	backoffFactor = 2
)

func RsyncFromServer(options RsyncOptions) (err error) {
	if options.User == "" {
		return errors.New("error: User is required")
	}
	if options.Host == "" {
		return errors.New("error: Host is required")
	}
	if options.DestinationDir == "" {
		options.DestinationDir = "temp_files"
	}

	// check if destination dir exists, if not create it
	if _, err := os.Stat(options.DestinationDir); os.IsNotExist(err) {
		fmt.Println("Creating destination directory: " + options.DestinationDir)
		err := os.MkdirAll(options.DestinationDir, 0755)
		if err != nil {
			log.Fatalf("Error creating destination directory: %v", err)
			return err
		}
	}

	var currentDelay time.Duration = initialDelay

	for retries := 0; retries < maxRetries; retries++ {
		if err := executeRsyncCommand(options); err != nil {
			log.Printf("Rsync attempt #%d failed: %v", retries+1, err)

			if retries < maxRetries-1 {
				log.Printf("Waiting for %v before next retry...", currentDelay)
				time.Sleep(currentDelay)
				currentDelay *= backoffFactor
			} else {
				return fmt.Errorf("rsync command failed after %d retries: %w", maxRetries, err)
			}
		} else {
			fmt.Println("✅ Rsync finished syncing the remote site directory to " + options.DestinationDir)
			return nil
		}
	}

	return fmt.Errorf("rsync command failed after reaching max retries")
}

func printOutput(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Print each line of the output
	}
}

func executeRsyncCommand(options RsyncOptions) error {
	rsyncCommand := "rsync"
	rsyncArgs := []string{
		"-azL", // archive, compress, and dereference symlinks
		"--progress",
		"-e", "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
		options.User + "@" + options.Host + ":" + os.Getenv("REMOTE_SITE_DIR"),
		options.DestinationDir,
	}
	cmd := exec.Command(rsyncCommand, rsyncArgs...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %w", err)
	}

	if options.Verbose {
		go printOutput(stdoutPipe)
	}
	go printOutput(stderrPipe)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting rsync command: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}

	return nil
}
