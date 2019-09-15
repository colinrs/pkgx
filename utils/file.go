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
