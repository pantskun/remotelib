package remotesftp

import (
	"path"
	"testing"
	"time"

	"github.com/pantskun/pathlib"
)

func TestUpload(t *testing.T) {
	config := SFTPConfig{
		Network:  "tcp",
		IP:       "192.168.62.11",
		Port:     "22",
		User:     "wx",
		Password: "1235",
		Timeout:  10 * time.Second,
	}

	interactor, err := NewInteractor(config)
	if err != nil {
		t.Fatal(err)
	}
	defer interactor.Close()

	err = interactor.Upload(path.Join(pathlib.GetModulePath("remotelib"), "remotesftp/test.txt"), "/home/wx/test")
	if err != nil {
		t.Log(err)
	}
}

func TestDownload(t *testing.T) {
	config := SFTPConfig{
		Network:  "tcp",
		IP:       "192.168.62.11",
		Port:     "22",
		User:     "wx",
		Password: "1235",
		Timeout:  10 * time.Second,
	}

	interactor, err := NewInteractor(config)
	if err != nil {
		t.Fatal(err)
	}
	defer interactor.Close()

	err = interactor.Download("/home/wx/test", pathlib.GetModulePath("remotelib"))
	if err != nil {
		t.Log(err)
	}
}
