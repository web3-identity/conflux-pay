package ginutils

import (
	"github.com/gin-gonic/gin"
)

func GetIdFromJwtClaim(c *gin.Context) uint {
	return c.GetUint("id")
}
