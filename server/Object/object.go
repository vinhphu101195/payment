package Object

import "time"

type PaymentMethod struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Order    string `json:"order"`
	ImgURL   string `json:"imgUrl"`
	Status   string `json:"status"`
	Platform string `json:"platform"`
	Note     string `json:"note"`
	Provider int    `json:"provider"`
}

type PaymentProvider struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PaymentAPI  string `json:"paymentAPI"`
	CallbackAPI string `json:"callbackAPI"`
	Note        string `json:"note"`
}

type PaymentItem struct {
	ID           int    `json:"id"`
	Method       string `json:"method"`
	Amount       int    `json:"amount"`
	Diamond      int    `json:"diamond"`
	DiamondBonus int    `json:"diamondBonus"`
	ImgURL       string `json:"imgUrl"`
	Status       string `json:"status"`
	Metadata     string `json:"metadata"`
}

type TransAction struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userId"`
	Amount        int       `json:"amount"`
	UserAmount    int       `json:"userAmount"`
	Currency      string    `json:"currency"`
	Diamond       int       `json:"diamond"`
	DiamondBonus  int       `json:"diamondBonus"`
	UserDiamond   int       `json:"userDiamond"`
	AppTransID    int       `json:"appTransId"`
	Source        string    `json:"source"`
	Status        string    `json:"status"`
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
	PaymentItemID int       `json:"paymentItemId"`
	SenderID      int       `json:"senderId"`
	Metadata      string    `json:"metadata"`
	Pin           string    `json:"pin"`
	Serie         string    `json:"serie"`
	ErrorMessage  string    `json:"errorMessage"`
	Provider      int       `json:"provider"`
}

func (pp PaymentProvider) TableName() string {
	return "payment_provider"
}

func (pm PaymentMethod) TableName() string {
	return "payment_provider"
}

func (pi PaymentItem) TableName() string {
	return "payment_provider"
}

func (t TransAction) TableName() string {
	return "transaction"
}
