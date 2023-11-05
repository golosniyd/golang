package main

import (
	"fmt"
	"log"
	"os/exec"
	"path"
)

// docker_deploy deploys a file to a Docker container.
//
// Parameters:
// - host: the host address of the Docker container.
// - filename: the name of the file to be deployed.
// - tmp_filepath: the temporary filepath on the host where the file will be copied.
// - target_filepath: the target filepath on the host where the file will be moved.
// - debug: a boolean indicating whether to enable debug mode.
//
// Returns: None.
func docker_deploy(host, filename, tmp_filepath, target_filepath string, debug bool) {
	if debug {
		log.Printf("Debug info for the func: docker_deploy\nHost: %s\nFilename: %s\nTmp_filepath: %s\nTarget_filepath: %s\n", host, filename, tmp_filepath, target_filepath)
	}
	log.Print("Deployment to a docker container has been started...")
	trimmed_filename := path.Base(filename)
	// fmt.Printf("Destination: \"%s\"", destination)
	err := copyFileWithSCP(host, filename, tmp_filepath)
	if err != nil {
		log.Fatalf("Failed to copy \"%s\" to directory on host: \"%s\" - \"%s\"\n", trimmed_filename, host, err)
	} else {
		log.Printf("Successfully copied \"%s\" to directory on host: \"%s\"\n", trimmed_filename, host)
	}
	err = moveFileWithSSH(host, tmp_filepath, target_filepath)
	if err != nil {
		log.Fatalf("Failed to deploy \"%s\" to docker tomcat volume on host: \"%s\" - \"%s\"\n", trimmed_filename, host, err)
	} else {
		log.Printf("Successfully deployed \"%s\" to docker tomcat volume on host: \"%s\"\n", trimmed_filename, host)
	}

}

// kubernetes_deploy deploys a file to a Kubernetes pod.
//
// Parameters:
// - appName: the name of the application.
// - filename: the name of the file to be deployed.
// - pathKuber: the path to the Kubernetes pod.
// - debug: a boolean flag indicating whether to enable debug mode.
func kubernetes_deploy(appName, filename, pathKuber string, debug bool) {
	if debug {
		log.Printf("Debug info for the func: kubernetes_deploy\nAppName: %s\nFilename: %s\nPathKuber: %s\n", appName, filename, pathKuber)
	}
	log.Print("Deployment to a kubernetes pod has been started...")
	pod_name, err := getPodName(appName)
	if err != nil || pod_name == "" {
		log.Fatal("Failed to get pod. Check if the pod is running.", err)
	}
	log.Printf("Pod was found: \"%s\"", pod_name)
	trimmed_filename := path.Base(filename)
	destination := fmt.Sprintf("%s:%s/%s", pod_name, pathKuber, trimmed_filename)
	log.Printf("Destination: \"%s\"", destination)
	err = deployFile(filename, destination)
	if err != nil {
		log.Fatalf("Failed to deploy \"%s\" to k8s pod: \"%s\" - \"%s\"\n", trimmed_filename, pod_name, err)
	} else {
		log.Printf("Successfully deployed \"%s\" to k8s pod: \"%s\"\n", trimmed_filename, pod_name)
	}
}

// deployFile copies a file from the source to the destination.
//
// Parameters:
// - src: the path of the source file.
// - dest: the path of the destination file.
//
// Returns:
// - error: an error if the file copy operation fails.
func deployFile(src, dest string) error {
	cmd := exec.Command("kubectl", "cp", src, dest)
	return cmd.Run()
}
