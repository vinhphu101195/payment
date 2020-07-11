package handle

import (
	"log"
	"payment/server/Object"
	"strconv"

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

//GetPaymentMethods ...
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

//GetPaymentMethodPopup for admin popup
func GetPaymentMethodPopup(ctx *gin.Context) {
	var pMethod []Object.PaymentMethod
	var pMethodPopup []Object.PaymentMethodPopup

	if err := db.Select([]string{"id,name"}).Where("status=\"active\"").Find(&pMethod).Scan(&pMethodPopup).Error; err != nil {

		ctx.JSON(200, gin.H{"error": 0, "data": pMethod})
	}
}

//GetPaymentItems get all items
func GetPaymentItems(ctx *gin.Context) {
	var pItem []Object.ShowPaymentItem
	// if err := db.Find(&pItem).Error; err != nil {
	// 	log.Println(err)
	// 	ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Item"}})
	// 	return
	// }

	db.Table("payment_item").Select("payment_item.id,payment_method.name as method_name, payment_item.method,payment_item.amount,payment_item.diamond,payment_item.diamond_bonus,payment_item.img_url,payment_item.status,payment_item.metadata").Joins("left join payment_method on payment_item.method = payment_method.id").Scan(&pItem)

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

//GetProviderPopup for admin popup
func GetProviderPopup(ctx *gin.Context) {
	var pProviders []Object.PaymentProvider
	var pProviderPopup []Object.PaymentProviderPopup

	if err := db.Select([]string{"id,name"}).Find(&pProviders).Scan(&pProviderPopup).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not find provider"}})
		return
	}
	if len(pProviderPopup) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No provider be found"}})
		return
	}
	ctx.JSON(200, gin.H{"error": 0, "data": pProviderPopup})
}

//GetTransaction ...
func GetTransaction(ctx *gin.Context) {
	var pTransaction []Object.TransAction
	const pagingSize = 10
	page, _ := strconv.Atoi(ctx.Param("page"))

	if err := db.Limit(pagingSize).Offset((page - 1) * pagingSize).Find(&pTransaction).Error; err != nil {
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
