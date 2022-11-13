package main

import (
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

func main() {
	client, err := ssh.Dial("tcp", "180.76.161.223:22", &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("lixiaohe@@52")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.TTY_OP_ISPEED: 14400,
	}

	if err = session.RequestPty("linux", 32, 160, modes); err != nil {
		log.Fatal(err)
	}

	if err = session.Shell(); err != nil {
		log.Fatal(err)
	}

	if err = session.Wait(); err != nil {
		log.Fatal(err)
	}
}
