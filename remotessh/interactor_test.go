package remotessh

import "testing"

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

	out, err := sshInteractor.Run("ls")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(out)
}
