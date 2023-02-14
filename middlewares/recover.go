package middlewares

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	rainbow_errors "github.com/web3-identity/conflux-pay/pay_errors"
	"github.com/web3-identity/conflux-pay/utils/ginutils"
)

func Recovery() gin.HandlerFunc {
	var buf bytes.Buffer
	return gin.CustomRecoveryWithWriter(&buf, gin.RecoveryFunc(func(c *gin.Context, err interface{}) {
		defer func() {
			fmt.Println(buf.String())
			logrus.WithField("recovered", buf.String()).Error("panic and recovery")
			buf.Reset()
		}()
		ginutils.RenderRespError(c, rainbow_errors.ERR_INTERNAL_SERVER_COMMON)
		c.Abort()
	}))
}
