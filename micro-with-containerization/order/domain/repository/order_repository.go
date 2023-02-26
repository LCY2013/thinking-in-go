package repository

import (
	"errors"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/model"
	"github.com/jinzhu/gorm"
)

type IOrderRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindOrderByID 根据OrderID查找Order信息
	FindOrderByID(int64) (*model.Order, error)
	// CreateOrder 创建Order
	CreateOrder(*model.Order) (int64, error)
	// DeleteOrderByID 根据OrderID删除Order
	DeleteOrderByID(int64) error
	// UpdateOrder 更新Order信息
	UpdateOrder(*model.Order) error
	// FindAll 查找所有Order信息
	FindAll() ([]model.Order, error)

	UpdateShipStatus(int64, int32) error
	UpdatePayStatus(int64, int32) error
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

type OrderRepository struct {
	db *gorm.DB
}

func (rep OrderRepository) InitTable() error {
	return rep.db.CreateTable(&model.Order{}, &model.OrderDetail{}).Error
}

func (rep OrderRepository) FindOrderByID(orderID int64) (order *model.Order, err error) {
	order = &model.Order{}
	return order, rep.db.Preload("OrderDetail").First(&order, orderID).Error
}

func (rep OrderRepository) CreateOrder(order *model.Order) (int64, error) {
	return order.ID, rep.db.Create(order).Error
}

func (rep OrderRepository) DeleteOrderByID(orderID int64) error {
	tx := rep.db.Begin()
	// 遇到错误就回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 彻底删除Order信息
	if err := tx.Unscoped().Where("id = ?", orderID).Delete(&model.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除OrderDetail信息
	if err := tx.Unscoped().Where("order_id = ?", orderID).Delete(&model.OrderDetail{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (rep OrderRepository) UpdateOrder(order *model.Order) error {
	return rep.db.Updates(order).Error
}

func (rep OrderRepository) FindAll() (orderAll []model.Order, err error) {
	return orderAll, rep.db.Preload("OrderDetail").Find(&orderAll).Error
}

// UpdateShipStatus 更新订单发货状态
func (rep OrderRepository) UpdateShipStatus(orderID int64, shipStatus int32) error {
	db := rep.db.Model(&model.Order{}).Where("id = ?", orderID).UpdateColumn("ship_status", shipStatus)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("更新订单发货状态失败")
	}

	return nil
}

// UpdatePayStatus 更新支付状态
func (rep OrderRepository) UpdatePayStatus(orderID int64, payStatus int32) error {
	db := rep.db.Model(&model.Order{}).Where("id = ?", orderID).UpdateColumn("pay_status", payStatus)
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("更新订单支付状态失败")
	}

	return nil
}
