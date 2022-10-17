package controllers

import "github.com/gin-gonic/gin"

type OrderController interface {
	MakeOrder(c *gin.Context)
	GetOrder(c *gin.Context)
}
