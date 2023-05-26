package model

type Order struct {
	OrderId       int64
	Weight        float64
	Regions       int64
	DeliveryHours []string
	Cost          int64
	CompletedTime *string
}
