package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")

	client := &http.Client{}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Setting up a channel to listen for interrupt signal (Ctrl + C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Using a channel to communicate when to stop the loop
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				runCloneBackupJob(client)
			case <-sigs:
				fmt.Println("\nReceived an interrupt, stopping...")
				done <- true
				return
			}
		}
	}()
	// Wait for signal to stop
	<-done
	fmt.Println("Program exiting")
}

func runCloneBackupJob(client *http.Client) {
	println("Running clone backup job...")

	sites, err := getSites(client)
	if err != nil {
		fmt.Println("Error getting sites:", err)
		return
	}

	// log the site name
	for _, site := range sites.Results {
		fmt.Println("Site:", site.Name)
	}

}

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

func getSites(client *http.Client) (SitesResponse, error) {
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
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	var response SitesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return SitesResponse{}, err
	}

	return response, nil

}
