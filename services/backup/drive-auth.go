package backupService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func initDriveService() (*drive.Service, error) {
	ctx := context.Background()

	// Read the OAuth 2.0 credentials file
	var b []byte
	var err error
	filePath := os.Getenv("GOOGLE_CLIENT_SECRET_JSON_FILE")
	for i := 0; i < 5; i++ { // Retry up to 5 times
		b, err = os.ReadFile(filePath)
		if err == nil {
			break // File read successfully
		}
		log.Printf("Attempt %d: Unable to read client secret file: %v", i+1, err)
		time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
	}
	if err != nil {
		log.Fatalf("Unable to read client secret file after several attempts: %v", err)
		return nil, err
	}

	// Get config from the JSON file
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	// Try to load saved token from file
	token, err := tokenFromFile("auth/token.json")
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

	return service, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Link the user to Google's consent page to ask for permission for the google drive scope.
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code from the code url param in the redirect url: \n%v\n", authURL)

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
	saveToken("auth/token.json", tok) // Save the token to a file
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	// make the auth directory if it doesn't exist
	if _, err := os.Stat("auth"); os.IsNotExist(err) {
		os.Mkdir("auth", 0755)
	}

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
