package utils

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

const (
	sshPass = "MyP@ss"
	sshUser = "eldius"
)

func Execute(ctx context.Context, host, port string) error {

	timeout := 10 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
	}

	fmt.Printf("timeout: %s...\n", timeout.String())

	c, err := ssh.Dial("tcp", net.JoinHostPort(host, port), &ssh.ClientConfig{
		User: "eldius",
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPass),
		},
		Timeout:         timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return fmt.Errorf("dial server: %w", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
		}
	}()

	//client, err := c.DialContext(ctx, "tcp", host+":"+port)
	//if err != nil {
	//	return fmt.Errorf("dial server: %w", err)
	//}
	//defer func() {
	//	if err := c.Close(); err != nil {
	//	}
	//}()
	//
	//_, err = client.Write([]byte("ls -lha /\n"))
	//if err != nil {
	//	return fmt.Errorf("write: %w", err)
	//}

	session, err := c.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	if err := session.Run("ls -lha ./"); err != nil {
		return fmt.Errorf("ls -lha ./: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
		}
	}()
	return nil
}
