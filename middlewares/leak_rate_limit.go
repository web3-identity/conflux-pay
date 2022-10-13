package middlewares

import (
	"flag"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// @see https://github.com/gin-gonic/examples/blob/master/ratelimiter/rate.go
var (
	limit ratelimit.Limiter
	rps   = flag.Int("rps", 100, "request per second")
)

func init() {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("[GIN] ")
	log.SetOutput(gin.DefaultWriter)
	limit = ratelimit.New(*rps)
}

func LeakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Print(color.CyanString("%v", now.Sub(prev)))
		prev = now
	}
}
