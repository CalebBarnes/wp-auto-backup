package backupService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
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