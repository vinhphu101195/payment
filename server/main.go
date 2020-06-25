package main

import (
	"payment/server/handle"

	"github.com/gin-gonic/gin"
)

func main() {
	handle.Init()
	r := gin.Default()
	r.GET("/", handle.GetPaymentMethod)

	r.Run(":8000")
}
