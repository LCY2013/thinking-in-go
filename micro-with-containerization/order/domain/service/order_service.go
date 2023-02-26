package service

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/repository"
)

type IOrderDataService interface {
	AddOrder(order *model.Order) (int64, error)
	DeleteOrder(orderID int64) error
	UpdateOrder(order *model.Order) error
	FindOrderByID(orderID int64) (order *model.Order, err error)
	FindAllOrder() ([]model.Order, error)

	UpdateShipStatus(orderID int64, shipStatus int32) error
	UpdatePayStatus(orderID int64, payStatus int32) error
}

// NewOrderDataService 创建用户数据服务
func NewOrderDataService(orderRepository repository.IOrderRepository) IOrderDataService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

type OrderService struct {
	orderRepository repository.IOrderRepository
}

func (ser OrderService) FindOrderByID(orderID int64) (order *model.Order, err error) {
	order = &model.Order{}
	return ser.orderRepository.FindOrderByID(orderID)
}

func (ser OrderService) AddOrder(order *model.Order) (int64, error) {
	// MQ
	return ser.orderRepository.CreateOrder(order)
}

func (ser OrderService) DeleteOrder(orderID int64) error {
	return ser.orderRepository.DeleteOrderByID(orderID)
}

func (ser OrderService) UpdateOrder(order *model.Order) error {
	// log
	return ser.orderRepository.UpdateOrder(order)
}

func (ser OrderService) FindAllOrder() ([]model.Order, error) {
	return ser.orderRepository.FindAll()
}

func (ser OrderService) UpdateShipStatus(orderID int64, shipStatus int32) error {
	return ser.orderRepository.UpdateShipStatus(orderID, shipStatus)
}

func (ser OrderService) UpdatePayStatus(orderID int64, payStatus int32) error {
	return ser.orderRepository.UpdatePayStatus(orderID, payStatus)
}
