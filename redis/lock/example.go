package lock

import (
	"fmt"
)

// Example ...
func Example() {
	var resLock *RdsLock
	var err error
	resLock, err = NewLock("127.0.0.1:6379")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	resLock.TryLock("resource", "token", 30)
	select {}
}
