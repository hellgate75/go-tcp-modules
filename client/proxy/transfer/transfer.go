package transfer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hellgate75/go-tcp-client/common"
	"github.com/hellgate75/go-tcp-common/log"
	commonnet "github.com/hellgate75/go-tcp-common/net"
	"io/ioutil"
	"os"
	"time"
)

type tranfer struct{
	logger log.Logger
}

var serverCommand string = "transfer-file"

func (tranfer *tranfer) SetLogger(logger log.Logger) {
	tranfer.logger = logger
}

func (tranfer *tranfer) SendMessage(conn *tls.Conn, params ...interface{}) error {
	var paramsLen int = len(params)
	if paramsLen < 2 {
		return errors.New(fmt.Sprintf("Insufficient number of parameters, expected 2 but give %v", paramsLen))
	}
	var origin string = fmt.Sprintf("%v", params[0])
	var destination string = fmt.Sprintf("%v", params[1])
	var perm = "0664"
	if len(params) > 2 {
		perm = fmt.Sprintf("%v", params[2])
	}
	var typeOfFile = "folder"

	var isMkdir bool = (origin == "folder")
	var data []byte
	if !isMkdir {
		info, err1 := os.Stat(origin)
		if err1 != nil {
			return err1
		}
		file, err2 := os.Open(origin)
		if err2 != nil {
			return err2
		}
		if !info.IsDir() {
			var err3 error
			data, err3 = ioutil.ReadAll(file)
			if err3 != nil {
				return err3
			}
			typeOfFile = "file"
		}
	}

	n0, err3b := commonnet.WriteString(serverCommand, conn)
	if err3b != nil {
		return err3b
	}
	if n0 == 0 {
		return errors.New(fmt.Sprintf("Unable to send command: %s", serverCommand))
	}
	n00, err3c := commonnet.WriteString(typeOfFile, conn)
	if err3c != nil {
		return err3c
	}
	if n00 == 0 {
		return errors.New(fmt.Sprintf("Unable to send type: %s", typeOfFile))
	}
	n1, err4 := common.WriteString(destination, conn)
	if err4 != nil {
		return err4
	}
	if n1 == 0 {
		return errors.New(fmt.Sprintf("Unable to send destination folder: %s", destination))
	}
	n2, err5 := commonnet.WriteString(perm, conn)
	if err5 != nil {
		return err5
	}
	if n2 == 0 {
		return errors.New(fmt.Sprintf("Unable to send file permissions: %s", perm))
	}
	if typeOfFile != "folder" {
		time.Sleep(3 * time.Second)
		n3, err6 := commonnet.Write(data, conn)
		if err6 != nil {
			return err6
		}
		if n3 == 0 {
			return errors.New(fmt.Sprintf("Unable to send data -> len: %v", len(data)))
		}
	}
	return nil
}
func (tranfer *tranfer) Helper() string {
	return "transfer-file [origin] [destination] [permissions]\n  Parameters:\n    [origin]           origin file path or folder path or 'folder' for mkdir\n    [destination]      remote file path or folder to create empty or copy recursively\n    [permissions]      (optional) remote file or folder permissions (default: 0664)\n"
}

func New() common.Sender {
	return &tranfer{}
}
