package remotessh

import (
	"testing"
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

	if err := sshInteractor.Run(cmds); err != nil {
		t.Fatal(err)
	}

	t.Log("stdout:", sshInteractor.GetStdout())
	t.Log("stderr:", sshInteractor.GetStderr())
}
