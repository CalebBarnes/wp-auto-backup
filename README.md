# WP Auto Backups

This Go-based tool schedules jobs to clone WordPress (WP) server backups to an external storage location. It's designed to work with any WP server that has the WP CLI installed, and it assumes that the server running this Go project has an SSH key added to the WP server for secure communication.

## Environment Variables

To configure the tool, set the following environment variables:

- `SITE_NAME` - The name of the site. This is used to name the backup files.
- `SSH_USER` - The SSH username for accessing the WP server.
- `SSH_HOST` - The hostname or IP address of the WP server.
- `SSH_KEY_PATH` - The path to the SSH key file for accessing the WP server. Defaults to `~/.ssh/id_rsa`.
- `REMOTE_SITE_DIR` - The directory path of the site on the remote WP server to backup.
- `GOOGLE_SERVICE_ACCOUNT_KEY_FILE` - The path to the Google Service Account key file for authentication with Google Drive.
- `GOOGLE_DRIVE_FOLDER_ID` - The ID of the Google Drive folder where backups will be stored.
- `VERBOSE` - (Optional) Set to `true` to enable verbose logging for troubleshooting and debugging.

### Example `.env.local` file for local development:

Create a `.env.local` file in your project's root directory with the following content:

```
SITE_NAME="your-site-name"
SSH_USER="your-ssh-user"
SSH_HOST="your-ssh-host"
SSH_KEY_PATH="~/.ssh/id_rsa"
REMOTE_SITE_DIR="/path/to/your/remote/site/directory"
GOOGLE_SERVICE_ACCOUNT_KEY_FILE="path/to/your-service-account-file.json"
GOOGLE_DRIVE_FOLDER_ID="your_google_drive_folder_id"
VERBOSE="true" # Optional for verbose logging
```

## Running Locally

To run the tool locally:

1. **Install Go**: Make sure you have Go installed on your local machine.
2. **Clone the Repository**: Clone this repository to your machine.
3. **Configure Environment Variables**: Create a `.env.local` file with the necessary environment variables.
4. **Run the Application**: Execute the app with `go run main.go`. Alternatively, for development with hot reloading, use `air`.

## Running in Production

To deploy the tool in a production environment:

1. **Build the Executable**: Compile the application into an executable using `go build`.
2. **Deploy**: Transfer the executable to your production server.
3. **Set Environment Variables**: Ensure all required environment variables are set in your server's environment.
4. **Run the Executable**: Start the application on your server.
