package service

import (
	"hot-coffee/models"
)

type OrderRepository interface {
	LoadOrders() ([]models.Order, error)
}

type MenuRepository interface {
	LoadMenuItems() ([]models.MenuItem, error)
}

type ReportService struct {
	orderRepo OrderRepository
	menuRepo  MenuRepository
}

func NewReportService(orderRepo OrderRepository, menuRepo MenuRepository) *ReportService {
	return &ReportService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *ReportService) GetTotalSales() (float64, error) {
	orders, err := s.orderRepo.LoadOrders()
	if err != nil {
		return 0, err
	}

	menuItems, err := s.menuRepo.LoadMenuItems()
	if err != nil {
		return 0, err
	}

	menuMap := make(map[string]float64)
	for _, item := range menuItems {
		menuMap[item.ID] = item.Price
	}

	var total float64
	for _, order := range orders {
		if order.Status != "closed" {
			continue
		}
		for _, item := range order.Items {
			price, ok := menuMap[item.ProductID]
			if !ok {
				continue
			}
			total += price * float64(item.Quantity)
		}
	}

	return total, nil
}

func (s *ReportService) GetPopularItems() ([]models.PopularItemReport, error) {
	orders, err := s.orderRepo.LoadOrders()
	if err != nil {
		return nil, err
	}

	menuItems, err := s.menuRepo.LoadMenuItems()
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, order := range orders {
		if order.Status != "closed" {
			continue
		}
		for _, item := range order.Items {
			counts[item.ProductID] += item.Quantity
		}
	}

	menuMap := make(map[string]string)
	for _, item := range menuItems {
		menuMap[item.ID] = item.Name
	}

	var result []models.PopularItemReport
	for id, count := range counts {
		result = append(result, models.PopularItemReport{
			ProductID: id,
			Name:      menuMap[id],
			Count:     count,
		})
	}

	return result, nil
}
