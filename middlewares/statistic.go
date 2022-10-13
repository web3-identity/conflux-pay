package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Statistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// claims := jwt.ExtractClaims(c)
		// userId := uint(claims[AppUserIdKey].(float64))
		// models.IncreaseStatistic(userId, c.Request.Method, c.FullPath())
	}
}
