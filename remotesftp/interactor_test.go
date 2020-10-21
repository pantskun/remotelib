package remotesftp

import (
	"testing"

	"github.com/pantskun/pathlib"
)

func TestInteractor(t *testing.T) {
	config := SFTPConfig{
		Network:  "tcp",
		IP:       "192.168.62.11",
		Port:     "22",
		User:     "wx",
		Password: "1235",
	}

	interactor, err := NewInteractor(config)
	if err != nil {
		t.Fatal(err)
	}
	defer interactor.Close()

	modulePath := pathlib.GetModulePath("remotelib")
	src := modulePath
	des := "/home/wx/"

	t.Log("src:", src, " des", des)

	err = interactor.Upload(src, des)
	if err != nil {
		t.Fatal(err)
	}
}
