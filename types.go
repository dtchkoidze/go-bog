package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type PaymentClient interface {
	Pay() *any
}

type BogPaymentClient struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	token        *string
}

func NewBogPaymentClient(clientID, clientSecret string) *BogPaymentClient {
	return &BogPaymentClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

type BogPaymentPayload struct {
	CallbackURL   string        `json:"callback_url"`
	LocalOrderId  int           `json:"external_order_id"`
	PurchaseUnits PurchaseUnits `json:"purchase_units"`
	RedirectURLs  RedirectURLs  `json:"redirect_urls"`
}

func NewBogPaymentPayload(cbURL string, loID int, pus PurchaseUnits, rURLS RedirectURLs) *BogPaymentPayload {
	return &BogPaymentPayload{
		LocalOrderId:  loID,
		PurchaseUnits: pus,
		RedirectURLs:  rURLS,
	}
}

func NewPaymentItem(q int, up decimal.Decimal, pID int) *PaymentItem {
	return &PaymentItem{
		Quantity:  q,
		UnitPrice: up,
		ProductID: pID,
	}
}

type BogPaymentResponse struct {
	ID    string          `json:"id"`
	Links BogPaymentLinks `json:"_links"`
}

type BogPaymentLinks struct {
	Details  Link `json:"details"`
	Redirect Link `json:"redirect"`
}

type Link struct {
	Href string `json:"href"`
}

type PurchaseUnits struct {
	Currency    string        `json:"currency"`
	TotalAmount int           `json:"total_amount"`
	Basket      []PaymentItem `json:"basket"`
}

type RedirectURLs struct {
	Success string `json:"success"`
	Fail    string `json:"fail"`
}

type PaymentItem struct {
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	ProductID int             `json:"product_id"`
}

type BogPaymentInfo struct {
	OrderID         string    `json:"order_id"`
	Industry        string    `json:"industry"`
	Capture         string    `json:"capture"`
	ExternalOrderID string    `json:"external_order_id"`
	Client          Client    `json:"client"`
	ZonedCreateDate time.Time `json:"zoned_create_date"`
	ZonedExpireDate time.Time `json:"zoned_expire_date"`
	OrderStatus     KV        `json:"order_status"`
	Buyer           Buyer     `json:"buyer"`
	PurchaseUnits   Purchase  `json:"purchase_units"`
	RedirectLinks   Redirects `json:"redirect_links"`
	PaymentDetail   Payment   `json:"payment_detail"`
	Discount        Discount  `json:"discount"`
	Actions         []Action  `json:"actions"`
	Lang            string    `json:"lang"`
	RejectReason    *string   `json:"reject_reason"`
}

type Client struct {
	ID      string `json:"id"`
	BrandKA string `json:"brand_ka"`
	BrandEN string `json:"brand_en"`
	URL     string `json:"url"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Buyer struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type Purchase struct {
	RequestAmount  decimal.Decimal `json:"request_amount"`
	TransferAmount decimal.Decimal `json:"transfer_amount"`
	RefundAmount   decimal.Decimal `json:"refund_amount"`
	CurrencyCode   string          `json:"currency_code"`
	Items          []Item          `json:"items"`
}

type Item struct {
	ExternalItemID    string          `json:"external_item_id"`
	Description       string          `json:"description"`
	Quantity          decimal.Decimal `json:"quantity"`
	UnitPrice         decimal.Decimal `json:"unit_price"`
	UnitDiscountPrice decimal.Decimal `json:"unit_discount_price"`
	VAT               decimal.Decimal `json:"vat"`
	VATPercent        decimal.Decimal `json:"vat_percent"`
	TotalPrice        decimal.Decimal `json:"total_price"`
	PackageCode       string          `json:"package_code"`
	TIN               *string         `json:"tin"`
	PINFL             *string         `json:"pinfl"`
	ProductDiscountID string          `json:"product_discount_id"`
}

type Redirects struct {
	Success string `json:"success"`
	Fail    string `json:"fail"`
}

type Payment struct {
	TransferMethod     KV     `json:"transfer_method"`
	Code               string `json:"code"`
	CodeDescription    string `json:"code_description"`
	TransactionID      string `json:"transaction_id"`
	PayerIdentifier    string `json:"payer_identifier"`
	PaymentOption      string `json:"payment_option"`
	CardType           string `json:"card_type"`
	CardExpiryDate     string `json:"card_expiry_date"`
	RequestAccountTag  string `json:"request_account_tag"`
	TransferAccountTag string `json:"transfer_account_tag"`
	SavedCardType      string `json:"saved_card_type"`
	ParentOrderID      string `json:"parent_order_id"`
}

type Discount struct {
	BankDiscountAmount   decimal.Decimal `json:"bank_discount_amount"`
	BankDiscountDesc     string          `json:"bank_discount_desc"`
	DiscountedAmount     decimal.Decimal `json:"discounted_amount"`
	OriginalOrderAmount  decimal.Decimal `json:"original_order_amount"`
	SystemDiscountAmount decimal.Decimal `json:"system_discount_amount"`
	SystemDiscountDesc   string          `json:"system_discount_desc"`
}

type Action struct {
	ActionID        string          `json:"action_id"`
	RequestChannel  string          `json:"request_channel"`
	Action          string          `json:"action"`
	Status          string          `json:"status"`
	ZonedActionDate time.Time       `json:"zoned_action_date"`
	Amount          decimal.Decimal `json:"amount"`
}
