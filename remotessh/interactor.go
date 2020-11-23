package remotessh

import (
	"bytes"
	"io"

	"golang.org/x/crypto/ssh"
)

type Interactor interface {
	Close()
	Run(cmds []string) error
	GetStdout() string
	GetStderr() string
}

type interactor struct {
	client *ssh.Client

	stdoutBuf bytes.Buffer
	stderrBuf bytes.Buffer
}

var _ Interactor = (*interactor)(nil)

func NewInteractor(config SSHConfig) (Interactor, error) {
	clientConfig := NewClientConfig(config.User, config.Password, 10, nil)

	client, err := ssh.Dial(config.Network, config.IP+":"+config.Port, clientConfig)
	if err != nil {
		return nil, err
	}

	return &interactor{client: client}, nil
}

func (i *interactor) Close() {
	i.client.Close()
}

func (i *interactor) Run(cmds []string) error {
	session, err := i.client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	var (
		stdout io.Reader
		stderr io.Reader
		stdin  io.WriteCloser
	)

	if stdout, err = session.StdoutPipe(); err != nil {
		return err
	}

	if stderr, err = session.StderrPipe(); err != nil {
		return err
	}

	if stdin, err = session.StdinPipe(); err != nil {
		return err
	}

	if err = session.Shell(); err != nil {
		return err
	}

	for _, cmd := range cmds {
		_, err = stdin.Write([]byte(cmd + "\n"))
		if err != nil {
			return err
		}
	}

	_, err = stdin.Write([]byte("exit\n"))
	if err != nil {
		return err
	}

	_, err = io.Copy(&i.stdoutBuf, stdout)
	if err != nil {
		return err
	}

	_, err = io.Copy(&i.stderrBuf, stderr)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) GetStdout() string {
	return i.stdoutBuf.String()
}

func (i *interactor) GetStderr() string {
	return i.stderrBuf.String()
}
