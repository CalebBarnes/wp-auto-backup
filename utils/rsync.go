package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type RsyncOptions struct {
	User           string
	Host           string
	DestinationDir string
	Verbose        bool
}

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

	rsyncCommand := "rsync"
	rsyncArgs := []string{
		"-azL", // archive, compress, and dereference symlinks (copy the actual files instead of symlinks)
		"--progress",
		"-e", "ssh",
		options.User + "@" + options.Host + ":" + os.Getenv("REMOTE_SITE_DIR"),
		options.DestinationDir,
	}
	cmd := exec.Command(rsyncCommand, rsyncArgs...)

	// fmt.Println("Executing command:", cmd.String())

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating stdout pipe: %v", err)
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creating stderr pipe: %v", err)
		return err
	}

	fmt.Println("⚡️Connecting to " + options.Host + " with rsync.")
	// fmt.Println("Executing command:", cmd.String())
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting rsync command: %v", err)
		return err
	}

	if options.Verbose {
		go printOutput(stdoutPipe)
	}
	go printOutput(stderrPipe)

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Rsync command failed: %v", err)
		return err
	}
	fmt.Println("Rsync command completed")
	return nil
}

func printOutput(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Print each line of the output
	}
}
