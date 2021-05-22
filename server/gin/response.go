package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Responses interface {
	Response(*gin.Context, interface{}, error )
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

func (d *DefaultResponses) Response(ctx *gin.Context, result interface{}, err error){
	// 自己的实现的 Responses 接口可以在这里实现错误err处理
	if out, ok := result.(Output);ok{
		out.Output() // 在输出的时候自己做转换
	}
	ctx.JSON(http.StatusOK, response{
		Code: 0,
		Message: "",
		Data: result,
	})
}

func newDefaultResponses() *DefaultResponses{
	return &DefaultResponses{}
}


func(engine *Engine) SetResponses(responses Responses){
	engine.responses = responses
}