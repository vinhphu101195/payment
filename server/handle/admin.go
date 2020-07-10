package handle

import (
	"log"
	"payment/server/Object"

	"github.com/gin-gonic/gin"
)

//AddPaymentMethod ...
func AddPaymentMethod(ctx *gin.Context) {
	var pMethod Object.PaymentMethod
	if err := ctx.ShouldBindJSON(&pMethod); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid method"}})
		return
	}
	pMethod.ID = 0

	db.Save(&pMethod)

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"insert_id": pMethod.ID}})
}

//AddPaymentItem ...
func AddPaymentItem(ctx *gin.Context) {
	var pItem Object.PaymentItem
	if err := ctx.ShouldBindJSON(&pItem); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid method"}})
		return
	}
	pItem.ID = 0

	if err := db.Save(&pItem).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not insert to database"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"insert_id": pItem.ID}})
}

func GetPaymentMethods(ctx *gin.Context) {
	provider := ctx.Query("provider")
	var pMethod []Object.PaymentMethod
	if err := db.Where("provider=?", provider).Find(&pMethod).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not find method"}})
		return
	}

	if len(pMethod) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": pMethod})
}

//GetPaymentItems get all items
func GetPaymentItems(ctx *gin.Context) {
	method := ctx.Query("method")
	var pItem []Object.PaymentItem
	if err := db.Where("method=?", method).Find(&pItem).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Item"}})
		return
	}

	if len(pItem) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No payment item be found"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": pItem})
}

//GetProviders get all Providers
func GetProviders(ctx *gin.Context) {
	var pProviders []Object.PaymentProvider
	if err := db.Find(&pProviders).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Providers"}})
		return
	}

	if len(pProviders) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No Providers be found"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": pProviders})

}

//GetTransaction ...
func GetTransaction(ctx *gin.Context) {
	var pTransaction []Object.TransAction
	if err := db.Find(&pTransaction).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Transaction"}})
		return
	}

	if len(pTransaction) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No Transaction be found"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": pTransaction})
}

//UpdatePaymentMethod ...
func UpdatePaymentMethod(ctx *gin.Context) {
	var pMethod Object.PaymentMethod
	if err := ctx.ShouldBindJSON(&pMethod); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid method"}})
		return
	}

	if err := db.Save(&pMethod).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not insert to database"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"updated_id": pMethod.ID}})
}

//UpdatePaymentItem ...
func UpdatePaymentItem(ctx *gin.Context) {
	var pItem Object.PaymentItem
	if err := ctx.ShouldBindJSON(&pItem); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid method"}})
		return
	}

	if err := db.Save(&pItem).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not insert to database"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"insert_id": pItem.ID}})
}
