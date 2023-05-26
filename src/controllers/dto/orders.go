package dto

type CreateOrderRequest struct {
	Orders []CreateOrderDto `json:"orders" validate:"required"`
}

type CreateOrderDto struct {
	Weight        float64  `json:"weight" validate:"required"`
	Regions       int64    `json:"regions" validate:"required"`
	DeliveryHours []string `json:"delivery_hours" validate:"required"`
	Cost          int64    `json:"cost" validate:"required"`
}

type OrderDto struct {
	OrderId       int64    `json:"order_id" validate:"required"`
	Weight        float64  `json:"weight" validate:"required"`
	Regions       int64    `json:"regions" validate:"required"`
	DeliveryHours []string `json:"delivery_hours" validate:"required"`
	Cost          int64    `json:"cost" validate:"required"`
	CompletedTime *string  `json:"completed_time,omitempty"`
}

type CompleteOrderRequestDto struct {
	CompleteInfo []CompleteOrder `json:"complete_info" validate:"required"`
}

type CompleteOrder struct {
	CourierId    int64  `json:"courier_id" validate:"required"`
	OrderId      int64  `json:"order_id" validate:"required"`
	CompleteTime string `json:"complete_time" validate:"required"` //todo add validation
}
