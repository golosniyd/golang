package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// getPodName returns the name of a pod.
//
// It takes in the appName string parameter, which is the name of the application.
// It returns a string, which is the name of the pod, and an error if there was any.
func getPodName(appName string) (string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-l", "app="+appName, "-o", "custom-columns=:metadata.name", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// findFileByPattern finds a file by a given pattern for a specific service.
//
// It takes two parameters:
// - pattern: a string representing the file pattern to search for.
// - service: a string representing the service name.
//
// It returns two strings:
// - old_filename: the name of the old file found matching the pattern.
// - new_filename: the name of the new file with the service name appended.
func findFileByPattern(pattern, service string) (string, string) {

	new_filename := "target/" + service + ".war"

	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

	if len(matches) == 0 {
		log.Fatalf("No files found matching pattern: %s", pattern)
	}

	log.Printf("Files found matching pattern \"%s\": \"%s\"", pattern, matches[0])
	old_filename := matches[0]
	return old_filename, new_filename
}

// rename renames the file from old_filename to new_filename.
//
// Parameters:
// - old_filename: the name of the file to be renamed.
// - new_filename: the new name for the file.
//
// Returns:
// None.
func rename(old_filename, new_filename string) {
	err := os.Rename(old_filename, new_filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File was renamed to a new name: \"%s\"", new_filename)
}

// getServiceName returns the name of the service.
//
// It does not take any parameters.
// It returns a string.
func getServiceName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	service_name := filepath.Base(dir)
	trimmed_service_name := strings.Replace(service_name, "wa-", "", 1)
	return trimmed_service_name
}

// credentialHelper returns the username from the environment variable "USERNAME".
//
// No parameter.
// Return type: string.
func credentialHelper() string {
	username := os.Getenv("USERNAME")
	// keyPath, err := getRSAKeyPath()
	// if err != nil {
	// 	fmt.Println("RSA key not found, please create RSA key first")
	// }
	return username
}

// func getRSAKeyPath() (string, error) {
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		return "", err
// 	}

// 	keyPath := filepath.Join(homeDir, ".ssh", "id_rsa")
// 	return keyPath, nil
// }

// readPassword returns a byte slice and an error.
//
// None.
// ([]byte, error)
func readPassword() ([]byte, error) {
	fmt.Print("Enter password: ")
	// Turn off echoing
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error reading password: %s", err)
	}
	fmt.Println()
	return bytePassword, nil
}
