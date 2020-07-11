package handle

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"payment/server/Object"
	database "payment/server/database"
	"strconv"
	"strings"
	"time"

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

func Pay(ctx *gin.Context) {
	var body map[string]interface{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Invalid request"}})
		return
	}

	var pItem *Object.PaymentItem
	pItem = new(Object.PaymentItem)
	db.First(&pItem, fmt.Sprint(body["item_id"]))

	if pItem == nil {
		ctx.JSON(200, gin.H{"error": 404, "data": gin.H{"error": "Payment item not found"}})
		return
	}

	var provider Object.PaymentProvider
	db.Raw("Select pr.* from payment_provider pr inner join payment_method pm on pm.provider=pr.id where pm.id=?",
		fmt.Sprint(body["method_id"])).Scan(&provider)

	resultURL, err := BeginTransAction(ctx, pItem, provider, body)
	if err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Connect to provider failed"}})
		return
	}
	if resultURL == "" {
		ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"message": "Payment Success"}})
	} else {
		ctx.JSON(200, gin.H{"error": 0, "data": gin.H{"message": "Payment Hafl Success", "payURL": resultURL}})
	}

}

//BeginTransAction ...
func BeginTransAction(ctx *gin.Context, pItem *Object.PaymentItem, provider Object.PaymentProvider, body map[string]interface{}) (string, error) {

	var trans Object.TransAction
	trans.PaymentItemID = pItem.ID
	trans.Pin = fmt.Sprint(body["pin"])
	trans.Amount = pItem.Amount
	trans.CreateAt = time.Now()
	trans.UpdateAt = time.Now()
	trans.Provider = provider.ID
	trans.Serie = fmt.Sprint(body["serie"])
	trans.Source = fmt.Sprint(body["method_name"])
	trans.Status = "created"
	trans.CreateAt = time.Now()
	trans.UpdateAt = time.Now()
	trans.UserAmount = pItem.Amount
	trans.UserID, _ = strconv.Atoi(fmt.Sprint(body["user_id"]))
	trans.Diamond = pItem.Diamond
	trans.DiamondBonus = pItem.DiamondBonus

	if err := db.Save(&trans).Error; err != nil {
		return "", err
	}

	info := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(provider.Metadata), &info); err != nil {
		return "", err
	}

	info["payment_api"] = provider.PaymentAPI
	info["callback_api"] = provider.CallbackAPI

	info["serie"] = fmt.Sprint(body["serie"])
	info["pin"] = fmt.Sprint(body["pin"])

	var err error
	var resultURL string
	switch strings.ToLower(provider.Name) {
	case "napthengay":
		err = napTheNgay(info, trans)
	case "thuthere":
		err = thuTheRe(info, trans)
	// case "vnpay":
	// 	err = vnPay(info, trans)
	case "momo":
		resultURL, err = useMomo(ctx, info, trans)
	}
	if err != nil {
		trans.Status = "failed"
		db.Save(&trans)
		return "", err
	}
	return resultURL, nil
}

func napTheNgay(info map[string]string, trans Object.TransAction) error {
	switch strings.ToLower(trans.Source) {
	case "viettel":
		info["card_id"] = "1"
	case "vinaphone":
		info["card_id"] = "3"
	case "mobiphone":
		info["card_id"] = "2"
	case "zing":
		info["card_id"] = "4"
	case "fpt":
		info["card_id"] = "5"
	case "vtc":
		info["card_id"] = "6"
	}
	var plaintText = fmt.Sprintf("%s%s%d%d%d%s%s%s%s",
		info["merchant_id"], info["api_mail"], trans.ID, info["card_id"], trans.Amount, info["pin"], info["serie"], "md5", info["secret_key"])
	key := getMD5Hash(plaintText)
	url := fmt.Sprintf("%s?merchant_id=%s&card_id=%s&seri_field=%s&pin_field=%s&trans_id=%d&data_sign=%s&algo_mode=md5&api_email=%s&card_value=%d",
		info["payment_api"], info["merchant_id"], info["card_id"], info["serie"], info["pin"], trans.ID, key, info["api_mail"], trans.Amount)

	log.Println(url)
	res, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		trans.Status = "failed"
		db.Save(&trans)
		return err
	}

	response := make(map[string]interface{}, 0)
	if err = json.Unmarshal(resBody, &response); err != nil {
		return err
	}
	if response["code"] == 100 {
		trans.Status = "success"
		db.Save(&trans)
		return nil
	}
	return fmt.Errorf("%s", response["msg"])
}

func thuTheRe(info map[string]string, trans Object.TransAction) error {
	url := fmt.Sprintf("%s?id=%s&serial=%d&code=%s&cash=%s&type=%d&ghichu=",
		info["payment_api"], info["id"], info["serie"], info["pin"], trans.Amount, strings.ToLower(trans.Source))
	res, err := http.Post(url, "application/x-www-form-urlencoded", nil)

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := make(map[string]interface{}, 0)
	if err = json.Unmarshal(resBody, &response); err != nil {
		return err
	}

	if response["err"] == 0 {
		trans.Status = "success"
		db.Save(trans)
		return nil
	}
	return fmt.Errorf("%s", response["msg"])
}

// func vnPay(info map[string]string, trans Object.TransAction) error {
// 	vnp_Params:=map[string]interface{}
// 	vnp_Params["vnp_Version"] = "2";
//     vnp_Params["vnp_Command"] = "pay";
// 	vnp_Params["vnp_TmnCode"] = info["vnp_TmnCode"]
// 	vnp_Params["vnp_Locale"] = "vn";
//     vnp_Params["vnp_CurrCode"] = "VND";
//     vnp_Params["vnp_TxnRef"] = trans.ID;
//     vnp_Params["vnp_OrderInfo"] = info["vnp_OrderInfo"]
//     vnp_Params["vnp_OrderType"] = info["vnp_OrderType"]
//     vnp_Params["vnp_Amount"] = trans.Amount * 100;
//     vnp_Params["vnp_ReturnUrl"] = "http://localhost:3000/"
//     vnp_Params["vnp_IpAddr"] = info["vnp_IpAddr"];
//     vnp_Params["vnp_CreateDate"] = time.Now().Format("yyyymmddHHmmss");
//     vnp_Params["vnp_BankCode"] = info["vnp_bankCode"];
// 	vnp_Params
// 	type SortBy []Type

// 	func (a SortBy) Len() int           { return len(a) }
// 	func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// 	func (a SortBy) Less(i, j int) bool { return a[i] < a[j] }
// 	return fmt.Errorf(url)
// }

func useMomo(ctx *gin.Context, info map[string]string, trans Object.TransAction) (string, error) {
	const apiEndPoint = "https://test-payment.momo.vn/gw_payment/transactionProcessor"
	const returnURL = "http://localhost:8000/momo-result"
	const notifyURL = "http://localhost:8000/"
	const requestType = "captureMoMoWallet"

	info["requestType"] = requestType
	info["returnUrl"] = returnURL
	info["notifyUrl"] = notifyURL
	info["orderId"] = strconv.Itoa(trans.ID)
	info["requestId"] = strconv.Itoa(trans.ID)
	info["amount"] = strconv.Itoa(trans.Amount)
	info["orderInfo"] = "nap " + strconv.Itoa(trans.Amount)
	//create signature
	var plaintText = fmt.Sprintf("partnerCode=%s&accessKey=%s&requestId=%d&amount=%d&orderId=%d&orderInfo=%s&returnUrl=%s&notifyUrl=%s&extraData=%s",
		info["partnerCode"], info["accessKey"], trans.ID, trans.Amount, trans.ID, info["orderInfo"], returnURL, notifyURL, info["extraData"])
	var secrectKey = info["secretKey"]
	signature := getHMACSHA256(plaintText, fmt.Sprint(secrectKey))
	info["signature"] = signature

	jsonBody, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	delete(info, "secretKey")
	resp, err := http.Post(info["payment_api"], "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	respone := make(map[string]interface{}, 0)
	if err = json.Unmarshal(body, &respone); err != nil {
		return "", err
	}
	if fmt.Sprint(respone["errorCode"]) == "0" {
		return fmt.Sprint(respone["payUrl"]), nil
	}
	return "", fmt.Errorf(fmt.Sprint(respone["localMessage"]))
}

//MoMoReponse ...
func MoMoReponse(ctx *gin.Context) {

	var pProviders Object.PaymentProvider
	var trans Object.TransAction
	db.First(&trans, ctx.Query("orderId"))

	signatureRespone := ctx.Query("signature")
	req := make(map[string]string, 0)
	listRespone := []string{"partnerCode", "accessKey",
		"requestId", "amount", "orderId", "orderInfo",
		"orderType", "transId", "message", "localMessage",
		"responseTime", "errorCode", "payType", "extraData"}
	query := ""
	i := 0
	for _, val := range listRespone {
		req[val] = ctx.Query(val)
		if i != 0 {

			query += "&" + val + "=" + req[val]
		} else {
			query += val + "=" + req[val]
			i = 1
		}
	}
	// get secretKey from provider
	if err := db.Select([]string{"metadata"}).Where("name=?", "Momo").Find(&pProviders).Error; err != nil {
		log.Println(err)
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "Can not find provider"}})
		return
	}
	metaData := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(pProviders.Metadata), &metaData); err != nil {
		return
	}

	// create signature add database
	signature := getHMACSHA256(query, metaData["secretKey"])

	if signature != signatureRespone {
		ctx.JSON(200, gin.H{"error": 500, "data": gin.H{"error": "wrong signature"}})

	}

	if ctx.Query("errorCode") == "0" {
		trans.Status = "success"
		trans.AppTransID, _ = strconv.Atoi(ctx.Query("transId"))
		log.Println(trans.AppTransID)
		db.Save(&trans)
		ctx.Redirect(http.StatusMovedPermanently, "http://localhost:3000/?payment=success")
		return
	}
	trans.Status = "fail"
	trans.Metadata = "errorMessage: " + ctx.Query("message")
	db.Save(&trans)
	ctx.Redirect(http.StatusMovedPermanently, "http://localhost:3000/?payment=fail")
}

func getHMACSHA256(text string, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(text))
	expectedMac := hex.EncodeToString(mac.Sum(nil))
	return expectedMac
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
