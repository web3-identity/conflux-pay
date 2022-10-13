package cns_errors

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type RainbowError int

type RainbowErrorInfo struct {
	Message        string
	HttpStatusCode int
}

type RainbowErrorDetailInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var rainbowErrorInfos = make(map[RainbowError]RainbowErrorInfo)

func (r RainbowError) HttpStatusCode() int {
	return rainbowErrorInfos[r].HttpStatusCode
}

func (r RainbowError) Error() string {
	return rainbowErrorInfos[r].Message
}

func (r RainbowError) RenderJSON(c *gin.Context) {
	httpStatusCode := rainbowErrorInfos[r].HttpStatusCode
	c.JSON(httpStatusCode, r.ErrorResponse())
}

func (r RainbowError) AbortWithRenderJSON(c *gin.Context) {
	debug.PrintStack()
	c.Abort()
	r.RenderJSON(c)
}

func (r RainbowError) ErrorResponse() *RainbowErrorDetailInfo {
	return &RainbowErrorDetailInfo{
		Code:    int(r),
		Message: r.Error(),
	}
}
