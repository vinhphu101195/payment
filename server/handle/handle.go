package handle

import (
	"encoding/json"
	"log"
	"payment/server/Object"
	database "payment/server/database"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

var db *gorm.DB

func Init() {
	// enverr := godotenv.Load("../.env")
	// if enverr != nil {
	// 	log.Fatal("Error loading .env file in config")
	// }
	var err error
	db, err = database.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Object.PaymentItem{})
	db.AutoMigrate(&Object.PaymentMethod{})
	db.AutoMigrate(&Object.PaymentProvider{})
	db.AutoMigrate(&Object.TransAction{})
}

func GetPaymentMethod(ctx *gin.Context) {

	var pMethod []Object.PaymentMethod
	db.Find(&pMethod)

	if len(pMethod) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found"}})
		return
	}

	data, err := json.Marshal(pMethod)
	if err != nil {
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot parse data to json"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": data})
}

func GetPaymentItem(ctx *gin.Context) {

	pmName := ctx.Param("paymentMethod")
	if len(pmName) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid payment method"}})
		return
	}
	var pItem []Object.PaymentItem
	db.Find(&pMethod)

	if len(pMethod) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found"}})
		return
	}

	data, err := json.Marshal(pMethod)
	if err != nil {
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot parse data to json"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": data})
}
