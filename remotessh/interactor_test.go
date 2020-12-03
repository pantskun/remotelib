package remotessh

import (
	"log"
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
		"echo test",
	}

	if err := sshInteractor.Run(cmds); err != nil {
		t.Fatal(err)
	}

	log.Println(sshInteractor.GetStdout())
}
