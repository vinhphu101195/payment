package handle

import (
	"log"
	"payment/server/Object"

	"github.com/gin-gonic/gin"
)

func AddPaymentMethod(ctx *gin.Context) {
	var pMethod Object.PaymentMethod
	if ctx.BindJSON(&pMethod) != nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid method"}})
		return
	}
	pMethod.ID = 0

	db.Save(&pMethod)

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"insert_id": pMethod.ID}})
}

func AddPaymentItem(ctx *gin.Context) {
	var pMethod Object.PaymentMethod
	if ctx.BindJSON(&pMethod) != nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid method"}})
		return
	}
	pMethod.ID = 0

	if err := db.Save(&pMethod).Error; err != nil {
		log.Printf(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not insert to database"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"insert_id": pMethod.ID}})
}
