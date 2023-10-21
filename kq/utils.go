package kq

import (
	"fmt"

	"github.com/colinrs/pkgx/utils"
)

func kqRecover(cleanups ...func()) func() {

	return func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
		if p := recover(); p != nil {
			fmt.Printf("recove:%+v, stack:%s\n", p, utils.Stack())
		}
	}
}
