package main

import (
	"fmt"
	"payment/server/handle"

	_ "payment/server/database"

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
	r.GET("/payment-item/:paymentMethodId", handle.GetPaymentItem)
	r.POST("/payment", handle.Pay)

	r.Run(":8000")
}
