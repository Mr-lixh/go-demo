package main

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Hello world\n")
	})

	// Usage: ssh_demo 127.0.0.1 -p 2222
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		return password == "123456"
	})))
}
