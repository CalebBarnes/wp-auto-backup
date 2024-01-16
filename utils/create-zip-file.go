package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateZipFile(zipFileName string, sourceDir string) (string, error) {
	fmt.Println("\n🗜️ Creating ZIP file:", zipFileName)

	zipFile, err := os.Create(zipFileName)
	if err != nil {
		fmt.Println("Failed to create zip file:", err)
		return "", err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Function to add files to the zip
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory
		if path == sourceDir {
			return nil
		}
		// Create a zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// Ensure the header has the correct relative path
		header.Name, err = filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return err
		}
		// Write the header
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		// Write the file content if it's not a directory
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("🙈 Failed to add files to zip:", err)
		fmt.Println("")
		return "", err
	}
	fmt.Println("✅ ZIP file created successfully:", zipFileName)
	return zipFileName, nil
}
