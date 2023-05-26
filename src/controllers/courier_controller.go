package controllers

import (
	"Ya.SumSchool23/controllers/dto"
	cerrors "Ya.SumSchool23/controllers/errors"
	"Ya.SumSchool23/rate_limiter"
	"Ya.SumSchool23/services"
	"Ya.SumSchool23/services/service_data"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

const (
	getCouriers = handlerName(iota)
	getCourierById
	getCourierMetaById
	postCouriers
)

type CourierController struct {
	courierService *services.CourierService
	rateLimiters   map[handlerName]*rate_limiter.RateLimiter
}

func NewCourierController(s *services.CourierService) *CourierController {
	return &CourierController{
		courierService: s,
		rateLimiters: map[handlerName]*rate_limiter.RateLimiter{
			getCouriers:        rate_limiter.NewRateLimiter(),
			getCourierById:     rate_limiter.NewRateLimiter(),
			getCourierMetaById: rate_limiter.NewRateLimiter(),
			postCouriers:       rate_limiter.NewRateLimiter(),
		},
	}
}

func (c *CourierController) GetCouriers(ctx echo.Context) error {
	if !c.rateLimiters[getCouriers].RegisterCall() {
		return cerrors.TooManyRequests.New("get couriers method overloaded")
	}

	var limit int64
	limitStr := ctx.QueryParam("limit")
	if limitStr == "" {
		limit = 1
	} else {
		l, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return cerrors.BadRequest.Wrapf(err, "cannot parse query param 'limit', got '%s'", limitStr)
		}
		limit = l
	}

	var offset int64
	offsetStr := ctx.QueryParam("offset")
	if offsetStr == "" {
		offset = 0
	} else {
		off, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return cerrors.BadRequest.Wrapf(err, "cannot parse query param 'offset', got '%s'", offsetStr)
		}
		offset = off
	}

	couriers, err := c.courierService.GetCouriers(limit, offset)
	if err != nil {
		return err
	}

	couriersDto := make([]dto.CourierDto, len(couriers))

	for i := 0; i < len(couriers); i++ {
		couriersDto[i].CourierId = couriers[i].CourierId
		couriersDto[i].CourierType = couriers[i].CourierType
		couriersDto[i].Regions = couriers[i].Regions
		couriersDto[i].WorkingHours = couriers[i].WorkingHours
	}
	return ctx.JSON(http.StatusOK, couriersDto)
}

func (c *CourierController) GetCourierById(ctx echo.Context) error {
	if !c.rateLimiters[getCourierById].RegisterCall() {
		return cerrors.TooManyRequests.New("get courier by id method overloaded")
	}

	idStr := ctx.Param("courier_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return cerrors.BadRequest.Wrapf(err, "cannot parse path param 'courier_id', got '%s'", idStr)
	}

	courier, err := c.courierService.GetCourierById(id)
	if err != nil {
		return err
	}

	courierDto := dto.CourierDto{
		CourierId:    courier.CourierId,
		CourierType:  courier.CourierType,
		Regions:      courier.Regions,
		WorkingHours: courier.WorkingHours,
	}
	return ctx.JSON(http.StatusOK, courierDto)
}

func (c *CourierController) GetCourierMetaById(ctx echo.Context) error {
	if !c.rateLimiters[getCourierMetaById].RegisterCall() {
		return cerrors.TooManyRequests.New("get courier meta by id method overloaded")
	}

	idStr := ctx.Param("courier_id")
	startDateStr := ctx.QueryParam("startDate")
	endDateStr := ctx.QueryParam("endDate")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return cerrors.BadRequest.Wrapf(err, "cannot parse path param 'courier_id', got '%s'", idStr)
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return cerrors.BadRequest.Wrapf(err, "cannot parse query param 'startDate', got '%s'", startDateStr)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return cerrors.BadRequest.Wrapf(err, "cannot parse query param 'endDate', got '%s'", endDateStr)
	}

	return cerrors.NotImplemented.Newf(
		"Sorry I did not implement and complete this method because of not enough time. Params: %s, %s, %s",
		id, startDate, endDate)
}

func (c *CourierController) PostCouriers(ctx echo.Context) error {
	if !c.rateLimiters[postCouriers].RegisterCall() {
		return cerrors.TooManyRequests.New("post couriers method overloaded")
	}

	createCourierRequest := new(dto.CreateCourierRequest)
	if err := ctx.Bind(createCourierRequest); err != nil {
		return cerrors.BadRequest.Wrap(err, "cannot parse create courier request")
	}
	if err := ctx.Validate(createCourierRequest); err != nil {
		return cerrors.BadRequest.Wrap(err, "invalid create courier request")
	}

	data := make([]service_data.NewCourierData, len(createCourierRequest.Couriers))
	for i := 0; i < len(data); i++ {
		data[i].CourierType = createCourierRequest.Couriers[i].CourierType
		data[i].Regions = createCourierRequest.Couriers[i].Regions
		data[i].WorkingHours = createCourierRequest.Couriers[i].WorkingHours
	}

	createdCouriers, err := c.courierService.CreateCouriers(data)
	if err != nil {
		return err
	}
	response := dto.CreateCouriersResponse{}
	response.Couriers = make([]dto.CourierDto, len(createdCouriers))

	for i := 0; i < len(createdCouriers); i++ {
		response.Couriers[i].CourierId = createdCouriers[i].CourierId
		response.Couriers[i].CourierType = createdCouriers[i].CourierType
		response.Couriers[i].Regions = createdCouriers[i].Regions
		response.Couriers[i].WorkingHours = createdCouriers[i].WorkingHours
	}
	return ctx.JSON(http.StatusOK, response)
}
