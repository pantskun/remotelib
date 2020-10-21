package remotesftp

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pantskun/pathlib"
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

func (i *interactor) Upload(src string, des string) error {
	isDir, err := pathlib.IsDir(src)
	if err != nil {
		return err
	}

	if err := i.sftpClient.MkdirAll(des); err != nil {
		return err
	}

	if isDir {
		return i.uploadDirectory(src, des)
	} else {
		return i.uploadFile(src, des)
	}
}

func (i *interactor) Download(src string, des string) error {
	return nil
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
			err := i.uploadDirectory(newSrc, newDes)
			if err != nil {
				return err
			}
		} else {
			// 若为文件，上传文件
			err := i.uploadFile(newSrc, newDes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
