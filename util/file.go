package util

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

//GetBaseDir 获取项目根目录
func GetBaseDir() (string, error) {
	basePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	res, err := filepath.EvalSymlinks(filepath.Dir(basePath))
	if err != nil {
		return "", err
	}

	return res, nil
}

func SaveToLocal(localFileName string, uploadFile multipart.File) error {
	f, err := os.OpenFile(localFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	if _, err = io.Copy(f, uploadFile); err != nil {
		return nil
	}
	//将文件指针指回文件开头则可重新读取文件
	if _, err = uploadFile.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err = buf.Flush(); err != nil {
		return err
	}
	return nil
}
