package main

import (
	"fmt"
	"payment/server/handle"

	_ "payment/server/database"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.Default())

	r.GET("/payment-methods", handle.GetPaymentMethod)
	r.GET("/payment-item/:paymentMethodId", handle.GetPaymentItem)
	r.POST("/payment", handle.Pay)

	r.GET("/payment-items", handle.GetPaymentItems)
	r.GET("/payment-providers", handle.GetProviders)
	r.GET("/transaction/:page", handle.GetTransaction)
	r.POST("/payment-method", handle.AddPaymentMethod)
	r.POST("/payment-item", handle.AddPaymentItem)
	r.PUT("/payment-method", handle.UpdatePaymentMethod)
	r.PUT("/payment-item", handle.UpdatePaymentItem)

	r.Run(":8000")
}
