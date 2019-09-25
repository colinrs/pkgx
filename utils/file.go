package utils

import (
	"io"
	"os"
)

// WriteMsgToFile ....
func WriteMsgToFile(filename string, msg string) (err error) {

	var file *os.File
	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, msg)
	if err != nil {
		return err
	}
	return file.Sync()
}

// GetFileSize ...
func GetFileSize(path string) (size int64, err error) {
	var fi os.FileInfo
	fi, err = os.Stat(path)
	if nil != err {
		return
	}

	size = fi.Size()
	return
}

// FileIsExist ...
func FileIsExist(path string) (isExist bool) {
	_, err := os.Stat(path)

	isExist = err == nil || os.IsExist(err)
	return
}
