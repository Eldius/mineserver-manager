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
		timeout = time.Until(deadline)
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
		_ = c.Close()
	}()

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
		_ = session.Close()
	}()
	return nil
}
