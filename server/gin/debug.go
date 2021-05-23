package gin

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	rgin "github.com/gin-gonic/gin"
)

const ginSupportMinGoVer = 14


// SetMode sets gin mode according to input string.
func SetMode(value string) {
	rgin.SetMode(value)
}


func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v <= ginSupportMinGoVer {
		debugPrint(`[WARNING] Now Gin requires Go 1.14+.`)
	}
	debugPrint(`[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.`)
}

func debugPrint(format string, values ...interface{}) {
	if rgin.IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(rgin.DefaultWriter, "[GIN-debug] "+format, values...)
	}
}

func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

`)
}


func debugPrintError(err error) {
	if err != nil {
		if rgin.IsDebugging() {
			fmt.Fprintf(rgin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
		}
	}
}
