package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

// Download ...
func Download(toFile, url string) error {
	out, err := os.Create(toFile)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.Body == nil {
		return fmt.Errorf("%s response body is nil", url)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// ReadBytes ...
func ReadBytes(cpath string) ([]byte, error) {
	if !IsFileExists(cpath) {
		return nil, fmt.Errorf("%s not exists", cpath)
	}

	if IsDirExists(cpath) {
		return nil, fmt.Errorf("%s not file", cpath)
	}

	return ioutil.ReadFile(cpath)
}

// ReadString ...
func ReadString(cpath string) (string, error) {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// ReadJSON ...
func ReadJSON(cpath string, cptr interface{}) error {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %s", cpath, err.Error())
	}

	err = json.Unmarshal(bs, cptr)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", cpath, err.Error())
	}

	return nil
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

// EnsureDir ...
func EnsureDir(path string) {

	if !FileIsExist(path) {
		os.MkdirAll(path, 0777)
	}
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
