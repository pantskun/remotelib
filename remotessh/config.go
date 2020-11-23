package remotessh

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Network  string
	IP       string
	Port     string
	User     string
	Password string
	Timeout  time.Duration
}

const defaultTimeout = 10

func defaultHostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func NewClientConfig(user string, password string, timeout time.Duration, publicKey ssh.PublicKey) *ssh.ClientConfig {
	if publicKey != nil {
		return &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.FixedHostKey(publicKey),
			Timeout:         defaultTimeout * time.Second,
		}
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: defaultHostKeyCallback,
		Timeout:         defaultTimeout * time.Second,
	}
}
