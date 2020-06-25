package main

import (
	"payment/server/handle"

	"github.com/gin-gonic/gin"
)

func main() {
	handle.Init()
	r := gin.Default()
	r.GET("/payment-method", handle.GetPaymentMethod)
	r.GET("/payment-item/:paymentMethod", handle.GetPaymentItem)

	r.Run(":8000")
}
