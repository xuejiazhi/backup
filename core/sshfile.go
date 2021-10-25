package core

import (
	"bibt-SpeedSkat/backup/utils"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
	"gopkg.in/ini.v1"
	"os"
	"path"
	"strings"
	"time"
)

var OldDict string

func SshCopy(item *ini.Section) {
	//set default
	user := item.Key("user").String()
	passwd := item.Key("passwd").String()
	host := item.Key("host").String()
	port, _ := item.Key("port").Int()
	port = cast.ToInt(utils.If(port == 0, 22, port))

	//conn
	sftpCon, err := sftpConnect(user, passwd, host, port)
	if err != nil {
		fmt.Println("Remote SSH Connect Is Error,", err.Error())
		return
	}

	OldDict = item.Key("dir").String()

	RangeRemoteDir(sftpCon, item.Key("dir").String())
}

func RangeRemoteDir(sftpCon *sftp.Client, dict string) {
	RemoteFiles, _ := sftpCon.ReadDir(dict)
	sshBar := utils.LocalBar{
		BarCount:    len(RemoteFiles),
		Start:       0,
		Notice:      "DownLoad "+dict + ":",
		Graph:       "#",
		NoticeColor: 2,
	}
	sshBar.GenBar()
	for _, backDir := range RemoteFiles {
		RemotePath := path.Join(dict, backDir.Name())
		if backDir.IsDir() {
			RangeRemoteDir(sftpCon, RemotePath)
		} else {
			path := strings.Replace(dict, "\\", "/", -1)
			destNewPath := strings.Replace(path, OldDict, BackTmpDir, -1)
			SshCopyFile(sftpCon, RemotePath, destNewPath)
		}
		sshBar.PrintBar()
	}
	sshBar.EndBar()
}

func SshCopyFile(sftpCon *sftp.Client, remotePath, localPath string) {
	remoteFile, err := sftpCon.Open(remotePath)
	defer remoteFile.Close()
	if err != nil {
		return
	}

	//分割path目录
	destSplitPathDirs := strings.Split(localPath, "/")
	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index <= len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := PathExists(destSplitPath)
			if b == false {
				//fmt.Println("Create Dict:" + destSplitPath)
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	localFilename := path.Base(remotePath)
	dstFile, err := os.Create(path.Join(localPath, localFilename))
	if err != nil {
		fmt.Println("Local File Create Fail", err.Error())
		return
	}
	defer dstFile.Close()

	if _, err := remoteFile.WriteTo(dstFile); err != nil {
		fmt.Println("File Write Is Error ", err.Error())
		return
	}
}

func sftpConnect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}
