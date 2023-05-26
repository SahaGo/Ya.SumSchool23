package services

import (
	"Ya.SumSchool23/repositories"
	"Ya.SumSchool23/services/model"
	"Ya.SumSchool23/services/service_data"
)

type CourierService struct {
	courierRepository *repositories.CourierRepository
}

func NewCourierService(r *repositories.CourierRepository) *CourierService {
	return &CourierService{
		courierRepository: r,
	}
}

func (s *CourierService) GetCourierById(id int64) (*model.Courier, error) {
	return s.courierRepository.GetCourierById(id)
}

func (s *CourierService) GetCouriers(limit int64, offset int64) ([]*model.Courier, error) {
	return s.courierRepository.GetCouriers(limit, offset)
}

func (s *CourierService) CreateCouriers(data []service_data.NewCourierData) ([]*model.Courier, error) {
	result := make([]*model.Courier, 0)

	courierIds, err := s.courierRepository.CreateCouriers(data)
	if err != nil {
		return nil, err
	}
	for _, id := range courierIds {
		courier, err := s.courierRepository.GetCourierById(id)
		if err != nil {
			return nil, err
		}
		result = append(result, courier)
	}
	return result, nil
}
