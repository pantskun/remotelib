package remotesftp

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pantskun/remotelib/remotessh"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Interactor interface {
	Close()
	Upload(src string, des string) error
	Download(src string, des string) error
}

type interactor struct {
	sftpClient *sftp.Client
	sshClient  *ssh.Client
}

var _ Interactor = (*interactor)(nil)

// NewInteractor
// 获取与服务器进行sftp交互的Interator.
func NewInteractor(config SFTPConfig) (Interactor, error) {
	sshClientConfig := remotessh.NewClientConfig(config.User, config.Password, config.Timeout, nil)

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

// Upload upload src(local) to des(remote).
func (i *interactor) Upload(src string, des string) error {
	if err := i.sftpClient.MkdirAll(des); err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return i.uploadDirectory(src, des)
	} else {
		return i.uploadFile(src, des)
	}
}

// Download download src(remote) to des(local).
func (i *interactor) Download(src string, des string) error {
	if err := os.MkdirAll(des, os.ModeDir); err != nil {
		return err
	}

	info, err := i.sftpClient.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return i.downloadDirectory(src, des)
	} else {
		return i.downloadFile(src, des)
	}
}

func (i *interactor) uploadFile(src string, des string) error {
	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}

	newDes := path.Join(des, path.Base(src))

	desfile, err := i.sftpClient.OpenFile(newDes, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		return err
	}

	_, err = io.Copy(desfile, srcfile)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) uploadDirectory(src string, des string) error {
	// 创建远程目录
	newDes := path.Join(des, path.Base(src))
	if err := i.sftpClient.MkdirAll(newDes); err != nil {
		return err
	}

	// 获取src目录内容
	fileinfos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	// 遍历src目录内容
	for _, fileinfo := range fileinfos {
		newSrc := path.Join(src, fileinfo.Name())

		if fileinfo.IsDir() {
			// 若为目录，递归目录
			if err := i.uploadDirectory(newSrc, newDes); err != nil {
				return err
			}
		} else {
			// 若为文件，上传文件
			if err := i.uploadFile(newSrc, newDes); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *interactor) downloadFile(src string, des string) error {
	srcfile, err := i.sftpClient.OpenFile(src, os.O_RDONLY)
	if err != nil {
		return err
	}

	newDes := path.Join(des, path.Base(src))

	desfile, err := os.Create(newDes)
	if err != nil {
		return err
	}

	_, err = io.Copy(desfile, srcfile)
	if err != nil {
		return err
	}

	return nil
}

func (i *interactor) downloadDirectory(src string, des string) error {
	newDes := path.Join(des, path.Base(src))
	if err := os.MkdirAll(newDes, os.ModeDir); err != nil {
		return err
	}

	fileinfos, err := i.sftpClient.ReadDir(src)
	if err != nil {
		return err
	}

	for _, fileinfo := range fileinfos {
		newSrc := path.Join(src, fileinfo.Name())

		if fileinfo.IsDir() {
			if err := i.downloadDirectory(newSrc, newDes); err != nil {
				return err
			}
		} else {
			if err := i.downloadFile(newSrc, newDes); err != nil {
				return err
			}
		}
	}

	return nil
}
