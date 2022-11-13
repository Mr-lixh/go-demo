package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"testing"
)

func TestBasic1(t *testing.T) {
	client, err := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
		User:            "lixiaohe",
		Auth:            []ssh.AuthMethod{ssh.Password("028815")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		t.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run("ls"); err != nil {
		t.Fatal(err)
	}

	fmt.Println(b.String())
}

func TestShell(t *testing.T) {
	client, err := ssh.Dial("tcp", "180.76.161.223:22", &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("lixiaohe@@52")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		t.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
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

	if err = session.RequestPty("linux", 32, 16, modes); err != nil {
		t.Fatal(err)
	}

	if err = session.Shell(); err != nil {
		t.Fatal(err)
	}

	if err = session.Wait(); err != nil {
		t.Fatal(err)
	}
}
