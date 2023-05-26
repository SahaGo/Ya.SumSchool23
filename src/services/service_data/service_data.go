package service_data

type NewCourierData struct {
	CourierType  string
	Regions      []int64
	WorkingHours []string
}

type NewOrderData struct {
	Weight        float64
	Regions       int64
	DeliveryHours []string
	Cost          int64
}

type NewCompleteOrderData struct {
	CourierId    int64
	OrderId      int64
	CompleteTime string
}

type NewOrderAssignResponseData struct {
	Date     string
	Couriers []NewCouriersGroupOrdersData
}

type NewGroupOrdersData struct {
	GroupOrderId int64
	Orders       []NewOrderDtoData
}

type NewCouriersGroupOrdersData struct {
	CourierId int64
	Orders    []NewGroupOrdersData
}

type NewOrderDtoData struct {
	OrderId       int64
	Weight        float64
	Regions       int64
	DeliveryHours []string
	Cost          int64
	CompletedTime *string
}
