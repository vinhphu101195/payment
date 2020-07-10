package Object

import (
	"time"
)

type PaymentMethod struct {
	ID           int           `json:"method_id"`
	Name         string        `json:"name"`
	Order        int           `json:"order"`
	ImgURL       string        `json:"img_url"`
	Status       string        `json:"status"`
	Platform     string        `json:"platform"`
	Note         string        `json:"note"`
	Provider     int           `json:"provider"`
	PaymentItems []PaymentItem `gorm:"foreignkey:Method"`
}

//PaymentMethodPopup for popup
type PaymentMethodPopup struct {
	ID   string `json:"method_id"`
	Name string `json:"name"`
}

//ShowPaymentMethod for admin
type ShowPaymentMethod struct {
	ID           int    `json:"method_id"`
	Provider     int    `json:"provider"`
	ProviderName string `json:"provider_name"`
	Name         string `json:"name"`
	Order        int    `json:"order"`
	ImgURL       string `json:"img_url"`
	Status       string `json:"status"`
	Platform     string `json:"platform"`
	Note         string `json:"note"`
}

type PaymentProvider struct {
	ID             int             `json:"provider_id"`
	Name           string          `json:"name"`
	PaymentAPI     string          `json:"payment_api"`
	CallbackAPI    string          `json:"callback_api"`
	Metadata       string          `json:"metadata"`
	Note           string          `json:"note"`
	PaymentMethods []PaymentMethod `gorm:"foreignkey:Provider"`
}

//PaymentProviderPopup for popup
type PaymentProviderPopup struct {
	ID   int    `json:"provider_id"`
	Name string `json:"name"`
}

type PaymentItem struct {
	ID           int    `json:"item_id"`
	Method       int    `json:"method"`
	Amount       int    `json:"amount"`
	Diamond      int    `json:"diamond"`
	DiamondBonus int    `json:"diamond_bonus"`
	ImgURL       string `json:"img_url"`
	Status       string `json:"status"`
	Metadata     string `json:"metadata"`
}

//ShowPaymentItem for admin
type ShowPaymentItem struct {
	ID           int    `json:"item_id"`
	MethodName   string `json:"method_name"`
	Method       string `json:"method"`
	Amount       int    `json:"amount"`
	Diamond      int    `json:"diamond"`
	DiamondBonus int    `json:"diamond_bonus"`
	ImgURL       string `json:"img_url"`
	Status       string `json:"status"`
	Metadata     string `json:"metadata"`
}

type TransAction struct {
	ID            int       `json:"trans_id"`
	UserID        int       `json:"user_id"`
	Amount        int       `json:"amount"`
	UserAmount    int       `json:"user_amount"`
	Currency      string    `json:"currency"`
	Diamond       int       `json:"diamond"`
	DiamondBonus  int       `json:"diamond_bonus"`
	UserDiamond   int       `json:"user_diamond"`
	AppTransID    int       `json:"app_trans_id"`
	Source        string    `json:"source"`
	Status        string    `json:"status"`
	CreateAt      time.Time `json:"create_at", gorm:"default:current_timestamp"`
	UpdateAt      time.Time `json:"update_at"`
	PaymentItemID int       `json:"item_id"`
	SenderID      int       `json:"sender_id"`
	Metadata      string    `json:"metadata"`
	Pin           string    `json:"pin"`
	Serie         string    `json:"serie"`
	ErrorMessage  string    `json:"error"`
	Provider      int       `json:"provider"`
}

func (pp PaymentProvider) TableName() string {
	return "payment_provider"
}

func (pm PaymentMethod) TableName() string {
	return "payment_method"
}

func (pi PaymentItem) TableName() string {
	return "payment_item"
}

func (t TransAction) TableName() string {
	return "transaction"
}

type NapTheNgayRequest struct {
	MerchantID string `json:"merchant_id"`
	CardID     int    `json:"card_id"`
	Seri       string `json:"seri_field"`
	Pin        string `json:"pin_field"`
	TransID    int    `json:"trans_id"`
	DataSign   string `json:"data_sign"`
	AlgoMode   string `json:"algo_mode"`
	APIEmail   string `json:"api_email"`
	CardValue  int    `json:"card_value"`
}

type Momo struct {
}
