package handle

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"payment/server/Object"
	database "payment/server/database"

	"github.com/jinzhu/gorm"

	"io/ioutil"
	"net/http"

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
	if err := db.Find(&pMethod).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not find method"}})
		return

	}

	if len(pMethod) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found"}})
		return
	}

	data, err := json.Marshal(pMethod)
	if err != nil {
		log.Println(err)
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
<<<<<<< HEAD
	db.Find(&pItem)

	if len(pItem) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found"}})
=======
	if err := db.Where("method=?", pmName).Find(&pItem); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot find Item"}})
		return
	}

	if len(pItem) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No payment item be found"}})
>>>>>>> 5b491df36977f0c4c275e0a5fac0d1b40819ee87
		return
	}

	data, err := json.Marshal(pItem)
	if err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Cannot parse data to json"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": data})
}

type payBody struct {
	PaymentMethodId int    `json:"pmid"`
	PaymentItemId   int    `json:"piid"`
	UserId          int    `json:"uid"`
	Serie           string `json:"serie"`
	Pin             string `json:"pin"`
}

func Pay(ctx *gin.Context) {
	var body payBody
	if ctx.BindJSON(&body) != nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid request"}})
		return
	}

	var pItem *Object.PaymentItem
	db.First(pItem, body.PaymentItemId)

	if pItem == nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Payment item not found"}})
		return
	}

	var provider *Object.PaymentProvider
	db.Raw("Select pr from payment_provider pr inner join payment_method pm on pm.provider=pr.id where pr.id=?",
		body.PaymentMethodId).Scan(provider)
	if provider == nil {
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Provider not found"}})
		return
	}

	if err := BeginTransAction(pItem, provider, body); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Connect to provider failed"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"message": "Payment Success"}})
}

func BeginTransAction(pItem *Object.PaymentItem, provider *Object.PaymentProvider, body payBody) error {

	var trans Object.TransAction
	trans.PaymentItemID = pItem.ID
	trans.Pin = body.Pin
	trans.Amount = pItem.Amount
	trans.Provider = provider.ID
	trans.Serie = body.Serie
	trans.Source = body.PaymentMethodId
	trans.Status = "created"
	trans.UserAmount = pItem.Amount
	trans.UserID = body.UserId
	trans.Diamond = pItem.Diamond
	trans.DiamondBonus = pItem.DiamondBonus

	if err := db.Save(trans).Error; err != nil {
		return err
	}

	info := make(map[string]interface{}, 0)
	if err = json.Unmarshal(provider.Metadata, info); err != nil {
		return err
	}

	res, err := napTheNgay(info, trans)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(resBody, NtnResponse); err != nil {
		return err
	}
	if NtnResponse.Code == 100 {
		trans.Status = "success"
		db.Save(trans)
		return nil
	}
	trans.Status = "failed"
	db.Save(trans)
	return fmt.Errorf("failed")
}

func napTheNgay(info map[string]interface{}, trans TransAction) {
	var plaintText = fmt.Sprintf("%s%s%d%d%d%s%s%s%s",
		info["merchant_id"], info["api_mail"], trans.ID, trans.PaymentItemID, trans.Amount, body.Pin, body.Serie, "md5", info["secret_key"])
	key := getMD5Hash(plaintText)
	url := fmt.Sprintf("%s?merchant_id=%s&card_id=%d&seri_field=%s&pin_field=%s&trans_id=%d&data_sign=%s&algo_mode=md5&api_email=%s&card_value=%d",
		info["url"], info["merchant_id"], trans.PaymentItemID, body.Serie, body.Pin, trans.ID, key, info["api_mail"], trans.Amount)

	http.Post(url, "application/json", nil)
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
