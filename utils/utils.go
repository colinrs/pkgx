package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"regexp"
	"text/template"
	"time"

	"github.com/colinrs/pkgx/logger"
)

// InfoRender ...
type InfoRender struct {
	FiveMin string
	Day     string
}

// NewInfoRender ...
func NewInfoRender() (nir *InfoRender) {

	now := time.Now()
	nir = &InfoRender{
		FiveMin: now.Format(FIVEMIN1),
	}
	return
}

// StrRender ...
func StrRender(src string) (s string, err error) {

	// https://stackoverflow.com/questions/23466497/how-to-truncate-a-string-in-a-golang-template
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		// TODO
		"truncate": func(s string) string {
			s = "new"
			return s
		},
	}
	infoForTemplate := NewInfoRender()
	tmpl, err := template.New("cronx").Funcs(funcMap).Parse(src)
	if err != nil {
		logger.Error("parsing: %s,err:%s", src, err.Error())
		return
	}

	err = tmpl.Execute(buf, infoForTemplate)
	if err != nil {
		logger.Error("execution: %s", err)
		return
	}

	s = buf.String()
	return
}

// GetWDPath gets the work directory path.
func GetWDPath() string {
	wd := os.Getenv("GOPATH")
	if wd == "" {
		panic("GOPATH is not setted in env.")
	}
	return wd
}

// IsDirExists judges path is directory or not.
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

// IsFileExists judges path is file or not.
func IsFileExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()

}

//IsNum judges string is number or not.
func IsNum(a string) bool {
	reg, _ := regexp.Compile("^\\d+$")
	return reg.MatchString(a)
}

// GetBytes interface 转 byte
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Struct2Map struct to map
func Struct2Map(obj interface{}) (map[string]interface{}, error) {
	var ret map[string]interface{}
	jsonStr, _ := json.Marshal(obj)
	err := json.Unmarshal(jsonStr, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SetValue(source interface{}, dest interface{}) error {
	receiverValue := reflect.ValueOf(dest).Elem()
	if receiverValue.CanSet() {
		if source != nil {
			valBytes, err := json.Marshal(source)
			if err != nil {
				return errors.New("cache:val_marshal_err")
			}
			err = json.Unmarshal(valBytes, dest)
			if err != nil {
				return errors.New("cache:val_unmarshal_err")
			}
			return nil
		}
		receiverValue.Set(reflect.Zero(receiverValue.Type()))
		return nil
	}
	return errors.New("cache:receiver_not_pointer")
}
