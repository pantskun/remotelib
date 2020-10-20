package remotessh

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

type Interactor interface {
	Close()
	Run(cmd string) (string, error)
}

type interactor struct {
	client *ssh.Client
}

var _ Interactor = (*interactor)(nil)

func NewInteractor(config SSHConfig) (Interactor, error) {
	clientConfig := NewClientConfig(config.User, config.Password, nil)

	client, err := ssh.Dial(config.Network, config.IP+":"+config.Port, clientConfig)
	if err != nil {
		return nil, err
	}

	return &interactor{client: client}, nil
}

func (i *interactor) Close() {
	i.client.Close()
}

func (i *interactor) Run(cmd string) (string, error) {
	session, err := i.client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}
