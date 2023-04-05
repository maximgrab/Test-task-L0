package structs

import (
	"time"
)

type Order struct {
	OrderUid          string
	TrackNumber       string
	Entry             string
	Locale            string
	InternalSignature string
	CustomerId        string `validate:"required,min=4,max=4"`
	DeliveryService   string `validate:"required,min=5,max=5"`
	ShardKey          string
	SmId              int
	DateCreated       time.Time `format:"2006-01-02T06:22:19Z" validate:"required"`
	OofShard          string
	Delivery          Delivery
	Payment           Payment
	Items             []Item
}

type Payment struct {
	Transaction  string
	RequestId    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}
type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Item struct {
	ChrtId      string
	TrackNumber string
	Price       string
	Rid         string
	Name        string
	Sale        string
	Size        string
	TotalPrice  string
	NmId        string
	Brand       string
	Status      string
}
