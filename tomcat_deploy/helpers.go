package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// getPodName retrieves the name of a pod associated with a given appName.
//
// Parameters:
// - appName: the name of the application associated with the pod.
//
// Returns:
// - string: the name of the pod.
// - error: an error if the command execution fails.
func getPodName(appName string) (string, error) {
	cmd := exec.Command(
		"kubectl",
		"get",
		"pods",
		"-l",
		fmt.Sprintf("app=%s", appName),
		"-o",
		"custom-columns=:metadata.name",
		"--no-headers",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// findFileByPattern finds a file by pattern and returns its old and new names.
//
// It takes two parameters: pattern (string) which is the pattern to match files
// and service (string) which is the name of the service.
//
// It returns two strings: old_filename (string) which is the name of the
// matched file and new_filename (string) which is the new name of the file.
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

// rename renames a file from old_filename to new_filename.
//
// Parameters:
// - old_filename: the name of the file to be renamed.
// - new_filename: the new name for the file.
//
// Return type: None.
func rename(old_filename, new_filename string) {
	err := os.Rename(old_filename, new_filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File was renamed to a new name: \"%s\"", new_filename)
}

// getServiceName retrieves the name of the current service.
//
// It does this by getting the current working directory, extracting the base name of the directory,
// and removing the prefix "wa-" from the base name.
//
// Returns:
// - string: The name of the current service.
func getServiceName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	service_name := filepath.Base(dir)
	trimmed_service_name := strings.Replace(service_name, "wa-", "", 1)
	return trimmed_service_name
}

// credentialHelper returns the username from the environment variables.
//
// No parameters.
// Returns a string.
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

// readPassword reads a password from the user input.
//
// It prompts the user to enter a password and then reads the input from the standard input.
// The password is read without echoing, meaning the characters entered by the user are not displayed on the screen.
// The function returns the password as a byte slice and any error that occurred during the reading process.
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

// validateConfig validates the given Config struct.
//
// It checks if the ContextPath field is not empty and the Host field is empty,
// and logs a fatal error if true. It also checks if the Host field is one of the
// valid values: "auto", "alpha", "rc", "dev", and logs a fatal error if not.
//
// Parameters:
// - cfg: A pointer to the Config struct to be validated.
//
// Return type: None.
func validateConfig(cfg *Config) {
	if cfg.ContextPath != "" && cfg.Host == "" {
		log.Fatal("Context flag can be used only with host flag: -h. Usage: -c=./path/to/context.xml -h rc")
	}
	validHosts := map[string]bool{
		"auto": true, "alpha": true, "rc": true, "dev": true,
	}
	if !validHosts[cfg.Host] && cfg.Host != "" {
		log.Fatal("Host must be one of: auto, alpha, rc, dev")
	}
}

// defineBuilding checks the configuration and builds the project if necessary.
//
// cfg: A pointer to the Config struct.
// Returns: None.
func defineBuilding(cfg *Config) {
	if cfg.Build {
		if cfg.ContextPath != "" {
			log.Println("Skip building maven project because context flag is used")
		} else {
			build()
		}
	}
}

// setFlags sets the flags for the given Config struct.
//
// It takes a pointer to a Config struct as its parameter.
// The function does not return anything.
func setFlags(cfg *Config) {
	flag.StringVar(&cfg.Host, "h", "", "<Host> Deploy service to a docker tomcat volume on selected host (e.g. rc, dev, auto, alpha).")
	flag.StringVar(&cfg.ContextPath, "c", "", "<Context> Relative path to a context.xml file.")
	flag.BoolVar(&cfg.Build, "b", true, "<Build> Build maven project before deployment")
	flag.BoolVar(&cfg.Kubernetes, "k", false, "<K8s> Deploy service to a k8s pod. Pod will be selected automatically (default: false)")
	flag.BoolVar(&cfg.Debug, "d", false, "<Debug> Debug mode. Will display path, filenames, etc (default: false)")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("No flags provided. Displaying help message.")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

// getWarFileDetails returns the details of a war file.
//
// It finds a file by pattern and returns the old name, name, docker path,
// and temporary path of the file.
//
// Returns:
// - oldName: the old name of the file.
// - name: the name of the file.
// - dockerPath: the path of the file in the docker environment.
// - tempPath: the path of the file in the temporary directory.
func getWarFileDetails() (string, string, string, string) {
	oldName, name := findFileByPattern(pattern, getServiceName())
	cleanName := strings.Replace(name, "target/", "", 1)
	dockerPath := fmt.Sprintf("%s/%s", pathDocker, cleanName)
	tempPath := fmt.Sprintf("%s/%s", tmpPath, cleanName)
	return oldName, name, dockerPath, tempPath
}

// getContextFilePaths returns the Docker context file path and the temporary context file path.
//
// No parameters.
// Returns two strings: dockerCtx and tmpCtx.
func getContextFilePaths() (string, string) {
	dockerCtx := fmt.Sprintf("%s/%s", pathContext, contextName)
	tmpCtx := fmt.Sprintf("%s/%s", tmpPath, contextName)
	return dockerCtx, tmpCtx
}
