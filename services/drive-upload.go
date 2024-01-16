package backupService

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func initDriveService() (*drive.Service, error) {
	ctx := context.Background()

	// ******* START Authenticate using oauth2 *******

	// Read the OAuth 2.0 credentials file
	b, err := os.ReadFile(os.Getenv("GOOGLE_CLIENT_SECRET_JSON_FILE"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// Get config from the JSON file
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	// Try to load saved token from file
	token, err := tokenFromFile("token.json")
	if err != nil {
		token = getTokenFromWeb(config) // Get a new token if not available
	}

	client := config.Client(ctx, token)

	// create new drive service
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}

	// ******* END Authenticate using oauth2 *******

	// ******* START Authenticate using the service account *******
	// serviceAccountKeyFile := os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_FILE")
	// b, err := os.ReadFile(serviceAccountKeyFile)
	// if err != nil {
	// 	log.Fatalf("Unable to read service account key file: %v", err)
	// 	return nil, err
	// }

	// Authenticate using the service account
	// driveConfig, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	// if err != nil {
	// 	log.Fatalf("Unable to parse service account key file to config: %v", err)
	// 	return nil, err
	// }
	// client := driveConfig.Client(ctx)

	// service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Drive client: %v", err)
	// 	return nil, err
	// }
	// ******* END Authenticate using the service account *******

	return service, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Link the user to Google's consent page to ask for permission for the google drive scope.
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	_, err := fmt.Scan(&authCode)
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
		return nil
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
		return nil
	}
	saveToken("token.json", tok) // Save the token to a file
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Println("Saving token to file", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Retrieves a token from a file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

type UploadFileOptions struct {
	FolderId string
	Filepath string
}

func UploadFile(options UploadFileOptions) (*drive.File, error) {
	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("☁️ Uploading file to Google Drive...")
	}
	service, err := initDriveService()
	if err != nil {
		return nil, fmt.Errorf("unable to init drive service: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Println("📁 Folder ID:", options.FolderId)
		fmt.Println("📄 Filepath:", options.Filepath)
	}

	// Open the file
	localFile, err := os.Open(options.Filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer localFile.Close()

	// Detect the content type of the file
	contentType := mime.TypeByExtension(filepath.Ext(options.Filepath))
	if contentType == "" {
		// Default to plain text if type could not be detected
		contentType = "text/plain"
	}

	// Get the filename from the Filepath
	_, filename := filepath.Split(options.Filepath)

	// Create a file on Google Drive
	driveFile := &drive.File{
		Name:    filename,
		Parents: []string{options.FolderId},
	}
	uploadedFile, err := service.Files.Create(driveFile).Media(localFile, googleapi.ContentType(contentType)).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Printf("✅ File '%s' uploaded with ID: %s\n", filename, uploadedFile.Id)
	}
	return uploadedFile, nil
}

type UploadBufferOptions struct {
	FolderId string
	Filename string
	Buffer   *bytes.Buffer
}

func UploadBuffer(options UploadBufferOptions) error {
	service, err := initDriveService()
	if err != nil {
		return fmt.Errorf("unable to init drive service: %v", err)
	}

	// Detect the content type of the file
	contentType := mime.TypeByExtension(filepath.Ext(options.Filename))
	if contentType == "" {
		// Default to plain text if type could not be detected
		contentType = "text/plain"
	}

	// Create a file on Google Drive
	driveFile := &drive.File{
		Name:    options.Filename,
		Parents: []string{options.FolderId},
	}
	file, err := service.Files.Create(driveFile).Media(bytes.NewReader(options.Buffer.Bytes()), googleapi.ContentType(contentType)).Do()
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}

	if os.Getenv("VERBOSE") == "true" {
		fmt.Printf("Buffer content '%s' uploaded with ID: %s\n", options.Filename, file.Id)
	}
	return nil
}
