package main

import (
	"fmt"
	"payment/server/handle"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print(err)
	}
	handle.Init()
	r := gin.Default()
	r.GET("/payment-method", handle.GetPaymentMethod)
	r.GET("/payment-item/:paymentMethod", handle.GetPaymentItem)

	r.Run(":8000")
}
