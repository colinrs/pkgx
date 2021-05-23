package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Responses interface {
	Response(*gin.Context, interface{}, error, Decode)
}

type Output interface {
	Output()
}


type DefaultResponses struct {

}

type response struct {
	Code uint64 `json:"code"`
	Message string  `json:"message"`
	Data interface{} `json:"data"`

}

func (d *DefaultResponses) Response(ctx *gin.Context, result interface{}, err error, DecodeImp Decode){
	// 可以在这里实现错误err处理
	code, message := DecodeImp.DecodeErr(err)
	if out, ok := result.(Output);ok{
		out.Output() // 在输出的时候自己做转换
	}
	ctx.JSON(http.StatusOK, response{
		Code: uint64(code),
		Message: message,
		Data: result,
	})
}

func newDefaultResponses() *DefaultResponses{
	return &DefaultResponses{}
}


func(engine *Engine) SetResponses(responses Responses){
	engine.responses = responses
}