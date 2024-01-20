package backupService

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"

	"golang.org/x/crypto/ssh"
)

type BackupDatabaseOptions struct {
	User string
	Host string
	Port string
}

func BackupDatabase(options BackupDatabaseOptions, timestamp string) {
	fmt.Println("🗄️ Starting database backup...")

	// Expand the tilde to the home directory path
	homeDir, err := user.Current()
	if err != nil {
		log.Fatalf("unable to get current user home directory: %v", err)
	}
	keyPath := os.Getenv("SSH_KEY_PATH")
	if keyPath == "" {
		keyPath = "~/.ssh/id_rsa" // set default
	}
	if keyPath[:2] == "~/" {
		keyPath = homeDir.HomeDir + keyPath[1:] // replace "~" with real home dir
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: options.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: This is insecure; in production, use a proper HostKeyCallback
	}
	// Connect to the SSH server
	conn, err := ssh.Dial("tcp", options.Host+":"+options.Port, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	fmt.Println("🔐 Connected with SSH to " + options.Host + ":" + options.Port)

	cmd := "wp db export -" // outputs the sql dump to stdout
	sess, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %v", err)
	}
	defer sess.Close()

	var stdoutBuf bytes.Buffer
	sess.Stdout = &stdoutBuf
	err = sess.Run(cmd)
	if err != nil {
		log.Printf("failed to run command: %v\nError Type: %T\nError Details: %+v\n", err, err, err)
		return
	}

	fmt.Println("📤 Uploading database dump to Google Drive...")

	fileName := fmt.Sprintf("%s-database-dump-%s.sql", options.User, timestamp)
	_, err = UploadBufferInSiteFolder(UploadBufferOptions{
		FolderId: os.Getenv("GOOGLE_DRIVE_FOLDER_ID"),
		Filename: fileName,
		Buffer:   &stdoutBuf,
	})
	if err != nil {
		log.Fatalf("Unable to upload buffer content: %v", err)
	}

	fmt.Println("✅ Database dump file uploaded to Google Drive: " + fileName)
	fmt.Println("")
}
