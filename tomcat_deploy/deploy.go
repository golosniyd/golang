package main

import (
	"fmt"
	"log"
	"os/exec"
	"path"
)

// docker_deploy deploys a file to a docker container on a specified host.
//
// Parameters:
//   - host: the address of the host where the docker container is running.
//   - filename: the name of the file to be deployed.
//   - tmp_filepath: the temporary filepath on the host where the file will be copied to.
//   - docker_filepath: the filepath inside the docker container where the file will be moved to.
func docker_deploy(host, filename, tmp_filepath, docker_filepath string) {
	log.Print("Deployment to a docker container has been started...")
	trimmed_filename := path.Base(filename)
	err := copyFileWithSCP(host, filename, tmp_filepath)
	if err != nil {
		log.Fatalf("Failed to copy \"%s\" to directory on host: \"%s\" - \"%s\"\n", trimmed_filename, host, err)
	} else {
		log.Printf("Successfully copied \"%s\" to directory on host: \"%s\"\n", trimmed_filename, host)
	}
	err = moveFileWithSSH(host, tmp_filepath, docker_filepath)
	if err != nil {
		log.Fatalf("Failed to deploy \"%s\" to docker tomcat volume on host: \"%s\" - \"%s\"\n", trimmed_filename, host, err)
	} else {
		log.Printf("Successfully deployed \"%s\" to docker tomcat volume on host: \"%s\"\n", trimmed_filename, host)
	}

}

// kubernetes_deploy deploys a file to a Kubernetes pod.
//
// It takes in the following parameters:
// - appName: the name of the application
// - filename: the name of the file to be deployed
// - pathKuber: the path in the Kubernetes cluster where the file will be deployed
//
// This function does the following:
// 1. Retrieves the pod name using the getPodName function.
// 2. Checks if the pod is running.
// 3. Trims the filename using the path.Base function.
// 4. Constructs the destination path in the Kubernetes pod.
// 5. Deploys the file to the Kubernetes pod using the deployFile function.
//
// If any errors occur during the deployment process, this function will log the error and exit.
// If the deployment is successful, this function will log a success message.
func kubernetes_deploy(appName, filename, pathKuber string) {
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

// deployFile deploys a file from the source to the destination.
//
// Parameters:
//   - src: The source file path.
//   - dest: The destination file path.
//
// Returns:
//   - error: An error if the deployment fails.
func deployFile(src, dest string) error {
	cmd := exec.Command("kubectl", "cp", src, dest)
	return cmd.Run()
}
