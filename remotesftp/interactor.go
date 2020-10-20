package remotesftp

import (
	"io"
	"os"

	"github.com/pantskun/remotelib/remotessh"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Interactor interface {
	Close()
	WriteFIle(src string, des string) error
	ReadFile(des string) ([]byte, error)
}

type interactor struct {
	sftpClient *sftp.Client
	sshClient  *ssh.Client
}

var _ Interactor = (*interactor)(nil)

func NewInteractor(config SFTPConfig) (Interactor, error) {
	sshClientConfig := remotessh.NewClientConfig(config.User, config.Password, nil)

	sshClient, err := ssh.Dial(config.Network, config.IP+":"+config.Port, sshClientConfig)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return &interactor{sftpClient: sftpClient, sshClient: sshClient}, nil
}

func (i *interactor) Close() {
	i.sftpClient.Close()
	i.sshClient.Close()
}

func (i *interactor) WriteFIle(src string, des string) error {
	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}

	desfile, err := i.sftpClient.Create(des)
	if err != nil {
		return err
	}

	_, err = io.Copy(desfile, srcfile)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) ReadFile(des string) ([]byte, error) {
	return nil, nil
}
