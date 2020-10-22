package remotessh

import (
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Network  string
	IP       string
	Port     string
	User     string
	Password string
}

func NewClientConfig(user string, password string, publicKey ssh.PublicKey) *ssh.ClientConfig {
	if publicKey != nil {
		return &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.FixedHostKey(publicKey),
			Timeout:         10 * time.Second,
		}
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
}
