package backupService

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

type BackupDatabaseOptions struct {
	User string
	Host string
	Port string
}

func BackupDatabase(options BackupDatabaseOptions, timestamp string) {
	fmt.Println("🗄️ Starting database backup...")
	key, err := os.ReadFile("ssh/wpengine_rsa")
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

	fmt.Println("Connected with SSH to " + options.Host + ":" + options.Port)

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
		log.Fatalf("failed to run command: %v", err)
	}

	if err := UploadBuffer(UploadBufferOptions{
		FolderId: os.Getenv("GOOGLE_DRIVE_FOLDER_ID"),
		Filename: fmt.Sprintf("%s-database-dump-%s.sql", options.User, timestamp),
		Buffer:   &stdoutBuf,
	}); err != nil {
		log.Fatalf("Unable to upload buffer content: %v", err)
	}
}
