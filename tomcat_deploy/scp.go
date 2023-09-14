package main

import (
	"context"
	"fmt"
	"log"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// copyFileWithSCP copies a file from the source path to the destination path using SCP protocol.
//
// Parameters:
// - host: the hostname or IP address of the remote server.
// - src: the path of the source file.
// - dest: the path of the destination file.
//
// Returns:
// - error: an error if any occurred during the file copying process.
func copyFileWithSCP(host, src, dest string) error {
	config, _, err := sshConfigHelper()
	if err != nil {
		return err
	}
	client := scp.NewClient(host, config)
	err = client.Connect()
	if err != nil {
		log.Println("Couldn't establish a connection to the remote server ", err)
		return err
	} else {
		log.Printf("Connected to %s\n", host)
	}

	f, err := os.Open(src)
	if err != nil {
		log.Println("Error while opening file ", err)
	}

	defer client.Close()
	defer f.Close()

	err = client.CopyFromFile(context.Background(), *f, dest, "0644")
	if err != nil {
		log.Println("Error while copying file ", err)
	}

	return err
}

// moveFileWithSSH moves a file from the source path to the destination path on a remote host using SSH.
//
// Parameters:
// - host: the address of the remote host.
// - src: the path of the source file.
// - dest: the path of the destination file.
//
// Returns:
// - error: an error if the file move operation fails.
func moveFileWithSSH(host, src, dest string) error {
	config, password, err := sshConfigHelper()
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("echo \"%s\" | sudo -S mv %s %s", password, src, dest)
	err = session.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}

// sshConfigHelper generates a SSH client configuration, password, and error.
//
// It calls the credentialHelper function to get the username.
// It calls the readPassword function to get the passwordBytes.
// It converts the passwordBytes to a string.
// It calls the auth.PasswordKey function to generate a passwordKeyClientConfig.
// It returns the passwordKeyClientConfig, the password, and an error.
func sshConfigHelper() (*ssh.ClientConfig, string, error) {
	username := credentialHelper()
	passwordBytes, err := readPassword()
	if err != nil {
		log.Print(err)
	}
	password := string(passwordBytes)
	// privateKeyClientConfig, err := auth.PrivateKey(username, keyPath, ssh.InsecureIgnoreHostKey())
	// if err != nil {
	// 	return nil, nil, "", err
	// }
	passwordKeyClientConfig, err := auth.PasswordKey(username, password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return nil, "", err
	}
	return &passwordKeyClientConfig, password, nil
}
