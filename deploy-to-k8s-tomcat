package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	build()
	rename(findByPattern())
	deploy()
}

// build is a Go function that executes a Maven build by running the "mvn clean package" command.
//
// This function does not take any parameters.
// It does not return any values.
func build() {
	cmd := exec.Command("mvn", "clean", "package")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	handleError(err)

	log.Println("Maven build completed")
}

// findByPattern finds files matching a given pattern and returns the old filename, new filename, and trimmed service name.
//
// No parameters.
// Returns three strings.
func findByPattern() (string, string, string) {
	pattern := "target/*.war"

	dir, err := os.Getwd()
	handleError(err)

	service_name := filepath.Base(dir)
	trim_service_name := strings.Replace(service_name, "wa-", "", 1)

	new_filename := "target/" + trim_service_name + ".war"

	matches, err := filepath.Glob(pattern)
	handleError(err)

	if len(matches) == 0 {
		log.Fatalf("No files found matching pattern: %s", pattern)
	}

	log.Printf("Files found matching pattern \"%s\": \"%s\"", pattern, matches[0])
	old_filename := matches[0]
	return old_filename, new_filename, trim_service_name
}

// rename renames a file from the old filename to the new filename.
//
// It takes three parameters:
// - old_filename: the name of the file to be renamed.
// - new_filename: the new name for the file.
// - trim_service_name: the service name to be trimmed.
//
// It does not return any value.
func rename(old_filename, new_filename, trim_service_name string) {
	err := os.Rename(old_filename, new_filename)
	handleError(err)

	log.Printf("Service name: \"%s\"", trim_service_name)
	log.Printf("File was renamed from: \"%s\" to a new name: \"%s\"", old_filename, new_filename)
}

// getPodName function returns the name of the pod for the given app name.
//
// It takes a string parameter 'appName' which represents the name of the app.
// It returns a string representing the name of the pod and an error if any.
func getPodName(appName string) (string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-l", "app="+appName, "-o", "custom-columns=:metadata.name", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// deploy function deploys a file to a Kubernetes pod running Tomcat.
//
// It searches for a file using a pattern, trims the filename, and deploys the file to the specified pod.
// The function returns an error if the deployment fails.
func deploy() {
	_, new_filename, _ := findByPattern()
	log.Print("Deployment has been started...")
	appName := "tomcat"

	pod_name, err := getPodName(appName)
	if err != nil || pod_name == "" {
		log.Fatal("Failed to get pod. Check if the pod is running.", err)
	}
	log.Printf("Pod was found: \"%s\"", pod_name)

	// trim filename to <service_name>.war
	trimmed_filename := path.Base(new_filename)

	destination := fmt.Sprintf("%s:/usr/local/tomcat/webapps/%s", pod_name, trimmed_filename)
	log.Printf("Destination: \"%s\"", destination)

	err = deployFile(new_filename, destination)
	if err != nil {
		log.Fatalf("Failed to deploy \"%s\" to \"%s\": \"%s\"\n", trimmed_filename, pod_name, err)
		os.Exit(1)
	} else {
		log.Printf("Successfully deployed \"%s\" to k8s pod: \"%s\"\n", trimmed_filename, pod_name)
	}
}

// deployFile function deploys a file from the source path to the destination path using kubectl.
//
// Parameters:
// - src: the path of the source file
// - dest: the path of the destination file
//
// Return type:
// - error: an error if there was a problem executing the kubectl command
func deployFile(src, dest string) error {
	cmd := exec.Command("kubectl", "cp", src, dest)
	return cmd.Run()
}

// handleError is a Go function that handles errors.
//
// It takes an error as a parameter and checks if it is not nil.
// If the error is not nil, it logs the error and exits the program.
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
