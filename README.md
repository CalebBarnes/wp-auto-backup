# WPE Backup Cloner

This tool, written in Go, schedules jobs that will clone WPE backups to an external storage location.

## Environment Variables

- `WPE_USER_ID`: The user ID for the WPE account. You can generate this in WPE on the API Access page.
- `WPE_PASSWORD`: The password for the WPE account. You can generate this in WPE on the API Access page.
- `WPE_INSTALLS`: A comma separated list of WPE installs to save backups for.
- `BACKUP_NOTIFICATION_EMAILS`: A comma separated list of emails to send backup notifications to.

Example `.env.local` file for local development:

```
WPE_USER_ID="3f8c8dd5-6b7d-4f7f-8a9f-9d5d2a8e4e51"
WPE_PASSWORD="7hVSkcmNZzgmKKRGEAXvNTvfxJxeX9zs"
WPE_INSTALLS="realcedar, realcedarstaging"
BACKUP_NOTIFICATION_EMAILS="you@example.com,anotherperson@example.com"
```

## Running Locally

1. Install Go
2. Clone this repo
3. Create a `.env.local` file with the environment variables listed above
4. Run `go run main.go` to run the app OR use air for hot reloading: `air`

## Running in Production

1. Build the app executable
2. Deploy the executable to your server
3. Add the environment variables listed above to your environment
4. Run the executable
