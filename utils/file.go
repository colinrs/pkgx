package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// WritePidFile ...
func WritePidFile(pidFile ...string) {
	fname := fmt.Sprintf("%s/pid", PWDDir())
	if len(pidFile) > 0 {
		fname = pidFile[0]
	}
	abs, err := filepath.Abs(fname)
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(abs)
	os.MkdirAll(dir, 0777)
	pid := os.Getpid()
	f, err := os.OpenFile(abs, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%d\n", pid))
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

// PWD gets compiled executable file absolute path.
func PWD() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// PWDDir gets compiled executable file directory.
func PWDDir() string {
	return filepath.Dir(PWD())
}
