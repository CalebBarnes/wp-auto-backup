package wpeService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Site struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Account   Account   `json:"account"`
	GroupName string    `json:"group_name"`
	Installs  []Install `json:"installs"`
}

type Account struct {
	ID string `json:"id"`
}

type Install struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	CName       string `json:"cname"`
	PHPVersion  string `json:"php_version"`
	IsMultisite bool   `json:"is_multisite"`
}

type SitesResponse struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
	Count    int    `json:"count"`
	Results  []Site `json:"results"`
}

func GetSites(client *http.Client) (SitesResponse, error) {
	userID := os.Getenv("WPE_USER_ID")
	password := os.Getenv("WPE_PASSWORD")

	req, err := http.NewRequest("GET", "https://api.wpengineapi.com/v1/sites", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return SitesResponse{}, err
	}

	req.SetBasicAuth(userID, password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return SitesResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return SitesResponse{}, err
	}

	var response SitesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return SitesResponse{}, err
	}

	return response, nil
}

type BackupRequest struct {
	Description        string   `json:"description"`
	NotificationEmails []string `json:"notification_emails"`
}

type BackupResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func CreateBackup(client *http.Client, install Install) (BackupResponse, error) {
	userID := os.Getenv("WPE_USER_ID")
	password := os.Getenv("WPE_PASSWORD")

	data := BackupRequest{
		Description: "Jambaree WPE backup cloner",
	}

	emails := strings.Split(os.Getenv("BACKUP_NOTIFICATION_EMAILS"), ",")
	for _, email := range emails {
		email = strings.TrimSpace(email) // Remove any leading or trailing spaces
		data.NotificationEmails = append(data.NotificationEmails, email)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return BackupResponse{}, err
	}

	endpoint := "https://api.wpengineapi.com/v1/installs/" + install.ID + "/backups"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return BackupResponse{}, err
	}
	req.SetBasicAuth(userID, password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return BackupResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return BackupResponse{}, err
	}

	var response BackupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return BackupResponse{}, err
	}

	fmt.Println("Backup ID:", response.ID)
	fmt.Println("Backup Status:", response.Status)
	return response, nil
}

func GetBackupStatus(client *http.Client, install Install, backupID string) (BackupResponse, error) {
	userID := os.Getenv("WPE_USER_ID")
	password := os.Getenv("WPE_PASSWORD")

	endpoint := "https://api.wpengineapi.com/v1/installs/" + install.ID + "/backups/" + backupID
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return BackupResponse{}, err
	}
	req.SetBasicAuth(userID, password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return BackupResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return BackupResponse{}, err
	}

	fmt.Println("Response body:", string(body))

	var response BackupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return BackupResponse{}, err
	}

	fmt.Println("Backup ID:", response.ID)
	fmt.Println("Backup Status:", response.Status)
	return response, nil
}
