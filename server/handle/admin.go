package handle

import (
	"log"
	"payment/server/Object"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type commonInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
	var pMethod []Object.PaymentMethod
	provider := []commonInfo{}
	db.Raw("select id,name from payment_provider").Scan(&provider)

	if err := db.Find(&pMethod).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not find method", "provider": provider}})
		return
	}

	if len(pMethod) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found", "provider": provider}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"method": pMethod, "provider": provider}})
}

//GetPaymentItems get all items
func GetPaymentItems(ctx *gin.Context) {
	var pItem []Object.PaymentItem

	method := []commonInfo{}
	db.Raw("select id,name from payment_method").Scan(&method)
	if err := db.Find(&pItem).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Item", "method": method}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"item": pItem, "method": method}})
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
	start := ctx.Query("startDay")
	end := ctx.Query("endDay")
	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		page = 1
	}

	if start == "" {
		start = "2020-01-01"
		end = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	} else if end == "" {
		end = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	var pTransaction []Object.TransAction
	if err := db.Limit(20).Offset((page-1)*20).Where(
		"create_at>? and create_at<?", start, end).Order("create_at desc").Find(&pTransaction).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Transaction"}})
		return
	}

	if len(pTransaction) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No Transaction be found"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"length": len(pTransaction), "transaction": pTransaction}})
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
