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

	r.GET("/payment-methods", handle.GetPaymentMethods)
	r.POST("/payment", handle.Pay)

	r.GET("/payment-items", handle.GetPaymentItems)
	r.GET("/payment-providers", handle.GetProviders)
	r.GET("/transaction", handle.GetTransaction)
	r.POST("/payment-method", handle.AddPaymentMethod)
	r.POST("/payment-item", handle.AddPaymentItem)
	r.PUT("/payment-method", handle.UpdatePaymentMethod)
	r.PUT("/payment-item", handle.UpdatePaymentItem)

	r.Run(":8000")
}
