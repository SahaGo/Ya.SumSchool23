package dto

type CreateCourierRequest struct {
	Couriers []CreateCourierDto `json:"couriers" validate:"required"`
}

type CreateCourierDto struct {
	CourierType  string   `json:"courier_type" validate:"required,oneof=FOOT BIKE AUTO"`
	Regions      []int64  `json:"regions" validate:"required"`
	WorkingHours []string `json:"working_hours" validate:"required,hh_mm_interval"`
}

type CreateCouriersResponse struct {
	Couriers []CourierDto `json:"couriers" validate:"required"`
}

type CourierDto struct {
	CourierId    int64    `json:"courier_id" validate:"required"`
	CourierType  string   `json:"courier_type" validate:"required,oneof=FOOT BIKE AUTO"`
	Regions      []int64  `json:"regions" validate:"required"`
	WorkingHours []string `json:"working_hours" validate:"required,hh_mm_interval"`
}

type GetCouriersResponse struct {
	Couriers []CourierDto `json:"couriers" validate:"required"`
	Limit    int32        `json:"limit" validate:"required"`
	Offset   int32        `json:"offset" validate:"required"`
}

type GetCourierMetaInfoResponse struct {
	CourierId    int64    `json:"courier_id" validate:"required"`
	CourierType  string   `json:"courier_type" validate:"required,oneof=FOOT BIKE AUTO"`
	Regions      []int32  `json:"regions" validate:"required"`
	WorkingHours []string `json:"working_hours" validate:"required,hh_mm_interval"`
	Rating       int32    `json:"rating"`
	Earnings     int32    `json:"earnings"`
}
