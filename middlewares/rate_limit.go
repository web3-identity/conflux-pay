package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var RateLimitMiddleware gin.HandlerFunc

func InitRateLimitMiddleware() {
	var rate = limiter.Rate{
		Period: 1 * time.Second,
		Limit:  viper.GetInt64("limits.rainbowApiPerSec"),
	}
	var store = memory.NewStore()
	var instance = limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
	RateLimitMiddleware = mgin.NewMiddleware(instance)

	logrus.WithField("limit config", rate).Info("set limit config")
}
