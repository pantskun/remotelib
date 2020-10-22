package remotessh

import (
	"bytes"
	"io"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestInteractor(t *testing.T) {
	sshConfig := SSHConfig{
		Network:  "tcp",
		IP:       "192.168.62.11",
		Port:     "22",
		User:     "wx",
		Password: "1235",
	}

	sshInteractor, err := NewInteractor(sshConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer sshInteractor.Close()

	cmds := []string{
		"cd /home/wx/CrawlerDemo",
		"go run ./start.go -n 4",
	}

	out, err := sshInteractor.Run(cmds)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(out)
}

func TestSSH(t *testing.T) {
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: "wx",
		Auth: []ssh.AuthMethod{
			ssh.Password("1235"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "192.168.62.11:22", config)
	if err != nil {
		t.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		t.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		t.Fatalf("unable to acquire stdout pipe: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		t.Fatalf("unable to acquire stdin pipe: %s", err)
	}

	// tm := ssh.TerminalModes{ssh.ECHO: 0}
	// if err = session.RequestPty("xterm", 80, 40, tm); err != nil {
	// 	t.Fatalf("req-pty failed: %s", err)
	// }

	if err := session.Shell(); err != nil {
		t.Fatal(err)
	}

	_, err = stdin.Write([]byte("echo $PATH & exit"))
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stdout); err != nil {
		t.Fatalf("reading failed: %s", err)
	}

	t.Log(buf.String())
}
