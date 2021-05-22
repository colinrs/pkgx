package gin

import (
	rgin "github.com/gin-gonic/gin"
)
func Logger() rgin.HandlerFunc{
	return rgin.Logger()
}

func Recovery() rgin.HandlerFunc{
	return rgin.Recovery()
}