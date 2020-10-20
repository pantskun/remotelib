package remotessh

import "golang.org/x/crypto/ssh"

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
		}
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}
