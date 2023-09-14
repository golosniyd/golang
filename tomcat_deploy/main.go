package main

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Host  string
	Build bool
}

// main is the entry point of the program.
//
// It parses command line flags to configure the application.
// The flags include:
// -h: Deploy to docker tomcat volume on selected host.
//
//	If not set, deploy will be performed to k8s pod.
//
// -b: Build maven project before deployment.
//
//	Default value is true.
//
// It initializes the necessary variables and constants.
// If the build flag is set to true, it calls the build function.
// It finds the file to be deployed using the specified pattern and service name.
// It trims the filename and sets the docker and temporary file paths.
// It renames the file.
// If the host is specified, it deploys to the docker tomcat volume using SSH.
// Otherwise, it deploys to the kubernetes pod.
//
// No parameters.
// No return types.
func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Host, "h", "", "Deploy to docker tomcat volume on selected host. If not set, deploy will be performed to k8s pod")
	flag.BoolVar(&cfg.Build, "b", true, "Build maven project before deployment")
	flag.Parse()
	host := cfg.Host + ":22"

	const pattern = "target/*.war"
	const appName = "tomcat"
	const pathDocker = "/var/lib/docker/volumes/docker_tomcat_webapps/_data"
	const pathKuber = "/usr/local/tomcat/webapps"
	const tmpPath = "/tmp"

	if cfg.Build {
		build()
	}
	old_filename, filename := findFileByPattern(pattern, getServiceName())
	trimmed_filename := strings.Replace(filename, "target/", "", 1)
	docker_filepath := fmt.Sprintf("%s/%s", pathDocker, trimmed_filename)
	tmp_filepath := fmt.Sprintf("%s/%s", tmpPath, trimmed_filename)
	rename(old_filename, filename)
	if cfg.Host != "" {
		docker_deploy(host, filename, tmp_filepath, docker_filepath)
	} else {
		kubernetes_deploy(appName, filename, pathKuber)
	}
}
