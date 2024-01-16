# WP Auto Backup

![wp-auto-backup-logo-256x256](https://github.com/CalebBarnes/wp-auto-backup/assets/24890515/2fa9242e-20f2-461b-be54-1b7138fc2840)

Go-powered tool for scheduling WordPress server & database backups to Google Drive.

WP Auto Backup is designed to work with any WP server that has the WP CLI installed and an SSH key added to the WP server for secure communication.

## Environment Variables

To configure the tool, set the following environment variables:

- `SITE_NAME` - The name of the site, used for naming the backup files and the Google Drive folder.
- `SSH_USER` - The SSH username for accessing the WP server.
- `SSH_HOST` - The hostname or IP address of the WP server.
- `SSH_KEY_PATH` - Optional. The path to the SSH key file for accessing the WP server. Defaults to `~/.ssh/id_rsa`.
- `REMOTE_SITE_DIR` - The directory path of the site on the remote WP server to backup.
- `GOOGLE_CLIENT_SECRET_JSON_FILE` - The path to the OAuth2 client secret JSON file for authentication with Google Drive.
- `GOOGLE_DRIVE_FOLDER_ID` - Optional. The ID of the Google Drive folder where backups folders will be stored. If not provided, individual backup folders are created in the root of Google Drive for each site.

### Example `.env.local` file for local development:

Create a `.env.local` file in your project's root directory, or copy the .env.local.example file in this repo and rename it to `.env.local`.

```
SITE_NAME="your-site-name"
SSH_USER="your-ssh-user"
SSH_HOST="your-ssh-host"
SSH_KEY_PATH="~/.ssh/id_rsa" # Optional
REMOTE_SITE_DIR="/path/to/your/remote/site/directory"
GOOGLE_CLIENT_SECRET_JSON_FILE="path/to/your/client_secret.json"
GOOGLE_DRIVE_FOLDER_ID="your_google_drive_folder_id"
VERBOSE="true" # Optional for verbose logging
BACKUP_ON_START="true" # Optional, creates a backup on start
```

## Running Locally

To run the tool locally:

1. **Install Go**: Ensure Go is installed on your local machine.
2. **Clone the Repository**: Clone this repository to your machine.
3. **Configure Environment Variables**: Create a `.env.local` file with the necessary environment variables.
4. **Run the Application**: Execute the app with `go run main.go`. Alternatively, for development with hot reloading, use `air`.

## Running in Production

To deploy the tool in a production environment:

1. **Build the Executable**: Compile the application into an executable using `go build`.
2. **Deploy**: Transfer the executable to your production server.
3. **Set Environment Variables**: Ensure all required environment variables are set in your server's environment.
4. **Run the Executable**: Start the application on your server.

## Running Locally

To run the tool locally:

1. **Install Go**: Ensure Go is installed on your local machine.
2. **Clone the Repository**: Clone this repository to your machine.
3. **Configure Environment Variables**: Create a `.env.local` file with the necessary environment variables.
4. **First-Time Setup**:
   - Run the application using `go run main.go`.
   - Upon the first run, you'll be prompted with a link in the console.
   - Open this link in your browser and grant access to your Google Drive.
   - Copy the authorization code provided by Google and paste it back into the console.
   - This process will generate and save an authentication token for subsequent runs.
5. **Run the Application**: After the initial setup, execute the app with `go run main.go`. For development with hot reloading, use `air`.

## Running in Production

To deploy the tool in a production environment:

1. **Build the Executable**: Compile the application into an executable using `go build`.
2. **Deploy**: Transfer the executable to your production server.
3. **Set Environment Variables**: Ensure all required environment variables are set in your server's environment.
4. **First-Time Setup on Production Server**:
   - Run the executable for the first time.
   - Follow the instructions in the console to open the provided OAuth link in a browser and grant access.
   - Copy the authorization code from Google and input it back into the server console.
   - This will authenticate the application with Google Drive and save the token for future use.
5. **Run the Executable**: Once set up, start the application on your server as needed.
