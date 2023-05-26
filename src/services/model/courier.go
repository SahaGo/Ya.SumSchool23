package model

type Courier struct {
	CourierId    int64
	CourierType  string
	Regions      []int64
	WorkingHours []string
}
