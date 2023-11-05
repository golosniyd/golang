package main

const pattern = "target/*.war"
const appName = "tomcat"
const pathDocker = "/var/lib/docker/volumes/docker_tomcat_webapps/_data"
const pathContext = "/var/lib/docker/volumes/docker_tomcat_conf/_data"
const pathKuber = "/usr/local/tomcat/webapps"
const tmpPath = "/tmp"
const contextName = "context.xml"

type Config struct {
	Host        string
	ContextPath string
	Build       bool
	Kubernetes  bool
	Debug       bool
}

// main is the entry point of the program.
//
// It initializes the configuration, sets the flags, validates the configuration,
// and performs the necessary actions based on the configuration.
// It does not take any parameters and does not return anything.
func main() {
	cfg := new(Config)
	setFlags(cfg)
	validateConfig(cfg)
	defineBuilding(cfg)

	host := cfg.Host + ".erp.sperasoft.com" + ":22"

	switch {
	case cfg.Host != "" && cfg.ContextPath != "":
		docker_context_filepath, tmp_context_filepath := getContextFilePaths()
		docker_deploy(host, cfg.ContextPath, tmp_context_filepath, docker_context_filepath, cfg.Debug)
	case cfg.Host != "":
		old_filename, filename, docker_filepath, tmp_filepath := getWarFileDetails()
		rename(old_filename, filename)
		docker_deploy(host, filename, tmp_filepath, docker_filepath, cfg.Debug)
	case cfg.Kubernetes:
		_, filename, _, _ := getWarFileDetails()
		kubernetes_deploy(appName, filename, pathKuber, cfg.Debug)
	}
}
