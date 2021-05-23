package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// Zip ...
func Zip(src, dst string) error {
	// 初始化给定的目录
	EnsureDir(src)
	EnsureDir(dst)
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			fmt.Printf("[Zip]关闭文件失败: %v", err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名(去除src)
		fh.Name = strings.TrimPrefix(path, src)

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return err
		}
		// 输出压缩的内容
		fmt.Printf("[Zip]成功压缩文件: %s, 共写入了 %d 个字符的数据\n", path, n)
		return nil
	})
}

