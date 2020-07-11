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

	r.GET("/", handle.Load)
	r.GET("/payment-methods", handle.GetPaymentMethods)
	r.POST("/payment", handle.Pay)

	r.GET("/payment-items", handle.GetPaymentItems)
	r.GET("/payment-providers", handle.GetProviders)
	r.GET("/transaction/:page", handle.GetTransaction)
	r.POST("/payment-method", handle.AddPaymentMethod)
	r.POST("/payment-item", handle.AddPaymentItem)
	r.PUT("/payment-method", handle.UpdatePaymentMethod)
	r.PUT("/payment-item", handle.UpdatePaymentItem)

	r.GET("/payment-method-popup", handle.GetPaymentMethodPopup)
	r.GET("/payment-provider-popup", handle.GetProviderPopup)

	r.GET("/momo-result", handle.MoMoReponse)

	r.Run(":8000")
}
