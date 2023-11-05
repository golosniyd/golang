package main

import (
	"log"
	"os"
	"os/exec"
)

func build() {
	log.Print("Building maven project...")
	cmd := exec.Command("mvn", "clean", "package")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Maven build completed")
}
