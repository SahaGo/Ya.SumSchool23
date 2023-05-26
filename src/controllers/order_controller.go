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
)

const (
	getOrders = handlerName(iota)
	getOrderById
	postOrders
	postOrdersComplete
)

type OrderController struct {
	orderService *services.OrderService
	rateLimiters map[handlerName]*rate_limiter.RateLimiter
}

func NewOrderController(s *services.OrderService) *OrderController {
	return &OrderController{
		orderService: s,
		rateLimiters: map[handlerName]*rate_limiter.RateLimiter{
			getOrders:          rate_limiter.NewRateLimiter(),
			getOrderById:       rate_limiter.NewRateLimiter(),
			postOrders:         rate_limiter.NewRateLimiter(),
			postOrdersComplete: rate_limiter.NewRateLimiter(),
		},
	}
}

func (c *OrderController) GetOrders(ctx echo.Context) error {
	if !c.rateLimiters[getOrders].RegisterCall() {
		return cerrors.TooManyRequests.New("get orders method overloaded")
	}

	limitStr := ctx.QueryParam("limit")
	var limit int64
	if limitStr == "" {
		limit = 1
	} else {
		l, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return cerrors.BadRequest.Wrapf(err, "cannot parse query param 'limit', got '%s'", limitStr)
		}
		limit = l
	}

	offsetStr := ctx.QueryParam("offset")
	var offset int64
	if offsetStr == "" {
		offset = 0
	} else {
		off, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return cerrors.BadRequest.Wrapf(err, "cannot parse query param 'offset', got '%s'", offsetStr)
		}
		offset = off
	}

	orders, err := c.orderService.GetOrders(limit, offset)
	if err != nil {
		return err
	}

	ordersDto := make([]dto.OrderDto, len(orders))
	for i := 0; i < len(orders); i++ {
		ordersDto[i].OrderId = orders[i].OrderId
		ordersDto[i].Weight = orders[i].Weight
		ordersDto[i].Regions = orders[i].Regions
		ordersDto[i].DeliveryHours = orders[i].DeliveryHours
		ordersDto[i].Cost = orders[i].Cost
		ordersDto[i].CompletedTime = orders[i].CompletedTime
	}
	return ctx.JSON(http.StatusOK, ordersDto)
}

func (c *OrderController) GetOrderById(ctx echo.Context) error {
	if !c.rateLimiters[getOrderById].RegisterCall() {
		return cerrors.TooManyRequests.New("get order by id method overloaded")
	}

	idStr := ctx.Param("order_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return cerrors.BadRequest.Wrapf(err, "cannot parse path param 'order_id', got '%s'", idStr)
	}

	order, err := c.orderService.GetOrderById(id)
	if err != nil {
		return err
	}

	orderDto := dto.OrderDto{
		OrderId:       order.OrderId,
		Weight:        order.Weight,
		Regions:       order.Regions,
		DeliveryHours: order.DeliveryHours,
		Cost:          order.Cost,
		CompletedTime: order.CompletedTime,
	}
	return ctx.JSON(http.StatusOK, orderDto)

}

func (c *OrderController) PostOrders(ctx echo.Context) error {
	if !c.rateLimiters[postOrders].RegisterCall() {
		return cerrors.TooManyRequests.New("post orders method overloaded")
	}

	createOrderRequest := new(dto.CreateOrderRequest)
	if err := ctx.Bind(createOrderRequest); err != nil {
		return cerrors.BadRequest.Wrap(err, "cannot parse create order request")
	}
	if err := ctx.Validate(createOrderRequest); err != nil {
		return cerrors.BadRequest.Wrap(err, "invalid create order request")
	}

	data := make([]service_data.NewOrderData, len(createOrderRequest.Orders))
	for i := 0; i < len(data); i++ {
		data[i].Weight = createOrderRequest.Orders[i].Weight
		data[i].Regions = createOrderRequest.Orders[i].Regions
		data[i].DeliveryHours = createOrderRequest.Orders[i].DeliveryHours
		data[i].Cost = createOrderRequest.Orders[i].Cost
	}

	createdOrders, err := c.orderService.CreateOrders(data)
	if err != nil {
		return err
	}

	response := make([]dto.OrderDto, len(createdOrders))
	for i := 0; i < len(createdOrders); i++ {
		response[i].OrderId = createdOrders[i].OrderId
		response[i].Weight = createdOrders[i].Weight
		response[i].Regions = createdOrders[i].Regions
		response[i].DeliveryHours = createdOrders[i].DeliveryHours
		response[i].Cost = createdOrders[i].Cost
		response[i].CompletedTime = createdOrders[i].CompletedTime
	}
	return ctx.JSON(http.StatusOK, response)

}

func (c *OrderController) PostOrdersComplete(ctx echo.Context) error {
	if !c.rateLimiters[postOrdersComplete].RegisterCall() {
		return cerrors.TooManyRequests.New("post orders complete method overloaded")
	}

	createCompleteOrderRequests := new(dto.CompleteOrderRequestDto)
	if err := ctx.Bind(createCompleteOrderRequests); err != nil {
		return cerrors.BadRequest.Wrap(err, "cannot parse create complete order request")
	}
	if err := ctx.Validate(createCompleteOrderRequests); err != nil {
		return cerrors.BadRequest.Wrap(err, "invalid create complete order requests")
	}

	data := make([]service_data.NewCompleteOrderData, len(createCompleteOrderRequests.CompleteInfo))
	for i := 0; i < len(data); i++ {
		data[i].CourierId = createCompleteOrderRequests.CompleteInfo[i].CourierId
		data[i].OrderId = createCompleteOrderRequests.CompleteInfo[i].OrderId
		data[i].CompleteTime = createCompleteOrderRequests.CompleteInfo[i].CompleteTime
	}

	createCompleteOrders, err := c.orderService.CreateCompleteOrder(data)
	if err != nil {
		return err
	}

	response := make([]dto.OrderDto, len(createCompleteOrders))
	for i := 0; i < len(createCompleteOrders); i++ {
		response[i].OrderId = createCompleteOrders[i].OrderId
		response[i].Weight = createCompleteOrders[i].Weight
		response[i].Regions = createCompleteOrders[i].Regions
		response[i].DeliveryHours = createCompleteOrders[i].DeliveryHours
		response[i].Cost = createCompleteOrders[i].Cost
		response[i].CompletedTime = createCompleteOrders[i].CompletedTime
	}
	return ctx.JSON(http.StatusOK, response)
}
