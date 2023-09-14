package main

import (
	"log"
	"os"
	"os/exec"
)

// build is a function that performs a Maven build by executing the "mvn clean package" command.
//
// This function does not take any parameters.
// It does not return any values.
func build() {
	cmd := exec.Command("mvn", "clean", "package")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Maven build completed")
}
