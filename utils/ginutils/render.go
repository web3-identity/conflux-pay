package ginutils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	rainbow_errors "github.com/web3-identity/conflux-pay/pay_errors"
)

func DataResponse(data interface{}) interface{} {
	return data
}

func ErrorResponse(code int, err error) gin.H {
	return gin.H{
		"code":    code,
		"message": err.Error(),
	}
}

func RenderResp(c *gin.Context, data interface{}, err error) {
	if err != nil {
		RenderRespError(c, err)
		return
	}
	RenderRespOK(c, data)
}

func RenderRespOK(c *gin.Context, data interface{}, httpStatusCode ...int) {
	statusCode := http.StatusOK

	if len(httpStatusCode) > 0 {
		statusCode = httpStatusCode[0]
	}
	c.JSON(statusCode, DataResponse(data))
}

// rainbowErrorCode 有值时，message 为 err message，如果 err 为 rainbow error, 则 status code 与 code 都来自 err, 否则来自rainbowErrorCode
// 否则 message 为 err message，status code 与 code 为 ERR_INTERNAL_SERVER_COMMON
func RenderRespError(c *gin.Context, err error, rainbowErrorCode ...rainbow_errors.RainbowError) {
	// logrus.WithField("error_stack", string(debug.Stack())).WithField("err", err).Info("render error")
	c.Error(err)
	c.Set("error_stack", fmt.Sprintf("%+v", errors.WithStack(err)))

	if re, ok := err.(rainbow_errors.RainbowError); ok {
		re.RenderJSON(c)
		return
	}

	_rainbowErrorCode := rainbow_errors.ERR_INTERNAL_SERVER_COMMON

	if len(rainbowErrorCode) > 0 {
		_rainbowErrorCode = rainbowErrorCode[0]
	}

	c.JSON(_rainbowErrorCode.HttpStatusCode(), ErrorResponse(int(_rainbowErrorCode), err))
}
