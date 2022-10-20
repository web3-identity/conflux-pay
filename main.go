package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wangdayong228/conflux-pay/config"
	_ "github.com/wangdayong228/conflux-pay/config"
	"github.com/wangdayong228/conflux-pay/logger"
	_ "github.com/wangdayong228/conflux-pay/logger"
	"github.com/wangdayong228/conflux-pay/middlewares"
	"github.com/wangdayong228/conflux-pay/models"
	"github.com/wangdayong228/conflux-pay/routers"
	"github.com/wangdayong228/conflux-pay/services"
	// "github.com/wangdayong228/conflux-pay/routers/assets"
	// "github.com/wangdayong228/conflux-pay/services"
)

func initGin() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middlewares.Logger())
	// engine.Use(gin.Recovery())
	engine.Use(middlewares.Recovery())
	return engine
}

// func init() {
// initConfig()
// logger.Init()
// middlewares.InitOpenJwtMiddleware()
// middlewares.InitRateLimitMiddleware()
// logrus.Info("init done")
// }

// @title       Rainbow-API
// @version     1.0
// @description The responses of the open api in swagger focus on the data field rather than the code and the message fields

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     127.0.0.1:8080
// @BasePath /v0
// @schemes  http https
func main() {
	config.Init()
	logger.Init()
	services.Init()

	models.ConnectDB()
	app := initGin()
	// var _ config.App
	// app.Use(middlewares.RateLimitMiddleware)
	routers.SetupRoutes(app)

	port := viper.GetString("port")
	if port == "" {
		logrus.Panic("port must be specified")
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)
	logrus.Info("Rainbow-API Start Listening and serving HTTP on ", address)
	err := app.Run(address)
	if err != nil {
		log.Panic(err)
	}
}
