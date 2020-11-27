package remotesftp

import "time"

type SFTPConfig struct {
	Network  string
	IP       string
	Port     string
	User     string
	Password string
	Timeout  time.Duration
}
