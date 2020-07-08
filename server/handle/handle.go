package handle

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"payment/server/Object"
	database "payment/server/database"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"io/ioutil"
	"net/http"
	"net/url"

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

func Load(ctx *gin.Context) {

	var pMethod []Object.PaymentMethod
	if err := db.Where("status=\"active\"").Find(&pMethod).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not find method"}})
		return
	}

	for i, _ := range pMethod {
		db.Model(&pMethod[i]).Association("PaymentItems").Find(&pMethod[i].PaymentItems)
	}

	if len(pMethod) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "No method be found"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": pMethod})
}

//GetPaymentItem ...
func GetPaymentItem(ctx *gin.Context) {
	pmID := ctx.Param("paymentMethodId")
	if len(pmID) == 0 {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid payment method"}})
		return
	}
	var pItem []Object.PaymentItem
	if err := db.Where("method=?", pmID).Find(&pItem).Error; err != nil {
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

type payBody struct {
	PaymentMethodId    int    `json:"method_id"`
	PaymentMethodName  string `json:"method_name"`
	PaymentMethodOrder int    `json:"method_order"`
	PaymentItemId      int    `json:"item_id"`
	UserId             int    `json:"user_id"`
	Serie              string `json:"serie"`
	Pin                string `json:"pin"`
}

func Pay(ctx *gin.Context) {
	var body payBody
	if ctx.BindJSON(&body) != nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid request"}})
		return
	}

	var pItem *Object.PaymentItem
	pItem = new(Object.PaymentItem)
	db.First(&pItem, body.PaymentItemId)

	if pItem == nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Payment item not found"}})
		return
	}

	var provider Object.PaymentProvider
	db.Raw("Select pr.* from payment_provider pr inner join payment_method pm on pm.provider=pr.id where pm.id=?",
		body.PaymentMethodId).Scan(&provider)
	// if provider == nil {
	// 	ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Provider not found"}})
	// 	return
	// }

	if err := BeginTransAction(pItem, provider, body); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Connect to provider failed"}})
		return
	}

	ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"message": "Payment Success"}})
}

func BeginTransAction(pItem *Object.PaymentItem, provider Object.PaymentProvider, body payBody) error {

	var trans Object.TransAction
	trans.PaymentItemID = pItem.ID
	trans.Pin = body.Pin
	trans.Amount = pItem.Amount
	trans.Provider = provider.ID
	trans.Serie = body.Serie
	trans.Source = body.PaymentMethodName
	trans.Status = "created"
	trans.CreateAt = time.Now()
	trans.UpdateAt = time.Now()
	trans.UserAmount = pItem.Amount
	trans.UserID = body.UserId
	trans.Diamond = pItem.Diamond
	trans.DiamondBonus = pItem.DiamondBonus

	if err := db.Create(&trans).Error; err != nil {
		return err
	}

	info := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(provider.Metadata), &info)

	info["serie"] = body.Serie
	info["pin"] = body.Pin

	var res *http.Response
	var err error
	switch strings.ToLower(provider.Name) {
	case "napthengay":
		res, err = napTheNgay(info, trans)
	case "thuthere":
		res, err = thuTheRe(info, trans)
	}
	if err != nil {
		trans.Status = "failed"
		db.Save(&trans)
		return err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		trans.Status = "failed"
		db.Save(&trans)
		return err
	}

	response := make(map[string]interface{}, 0)
	json.Unmarshal(resBody, &response)
	if response["code"] == 100 {
		trans.Status = "success"
		db.Save(trans)
		return nil
	}
	trans.Status = "failed"
	db.Save(&trans)
	return fmt.Errorf(fmt.Sprint(response["msg"]))
}

func napTheNgay(info map[string]interface{}, trans Object.TransAction) (*http.Response, error) {
	var plaintText = fmt.Sprintf("%s%s%d%d%d%s%s%s%s",
		info["merchant_id"], info["api_mail"], trans.ID, trans.PaymentItemID, trans.Amount, info["pin"], info["serie"], "md5", info["secret_key"])
	info["card_id"] = 1
	info["trans_id"] = trans.ID
	info["data_sign"] = getMD5Hash(plaintText)
	info["card_value"] = trans.Amount
	info["url"] = "http://api.napthengay.com/v2/"
	baseUrl, _ := url.Parse(fmt.Sprint(info["url"]))
	params := url.Values{}
	for key, val := range info {
		params.Add(key, fmt.Sprint(val))
	}

	baseUrl.RawQuery = params.Encode()
	return http.Post(baseUrl.String(), "application/x-www-form-urlencoded", nil)
}

func thuTheRe(info map[string]interface{}, trans Object.TransAction) (*http.Response, error) {
	return nil, fmt.Errorf("dasdad")
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
