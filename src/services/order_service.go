package services

import (
	"Ya.SumSchool23/repositories"
	"Ya.SumSchool23/services/model"
	"Ya.SumSchool23/services/service_data"
)

type OrderService struct {
	orderRepository *repositories.OrderRepository
}

func NewOrderService(r *repositories.OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: r,
	}
}

func (s *OrderService) GetOrderById(id int64) (*model.Order, error) {
	return s.orderRepository.GetOrderById(id)
}

func (s *OrderService) GetOrders(limit, offset int64) ([]*model.Order, error) {
	return s.orderRepository.GetOrders(limit, offset)
}

func (s *OrderService) CreateOrders(data []service_data.NewOrderData) ([]*model.Order, error) {
	result := make([]*model.Order, 0)

	orderIds, err := s.orderRepository.CreateOrders(data)
	if err != nil {
		return nil, err
	}
	for _, id := range orderIds {
		order, err := s.orderRepository.GetOrderById(id)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, nil
}

func (s *OrderService) CreateCompleteOrder(data []service_data.NewCompleteOrderData) ([]*model.Order, error) {
	result := make([]*model.Order, 0)

	orderIds, err := s.orderRepository.CreateCompleteOrder(data)
	if err != nil {
		return nil, err
	}
	for _, id := range orderIds {
		order, err := s.orderRepository.GetOrderById(id)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, nil
}
