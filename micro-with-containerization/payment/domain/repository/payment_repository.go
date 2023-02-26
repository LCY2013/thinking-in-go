package repository

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/model"
	"github.com/jinzhu/gorm"
)

type IPaymentRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindPaymentByID 根据PaymentID查找Payment信息
	FindPaymentByID(int64) (*model.Payment, error)
	// CreatePayment 创建Payment
	CreatePayment(*model.Payment) (int64, error)
	// DeletePaymentByID 根据PaymentID删除Payment
	DeletePaymentByID(int64) error
	// UpdatePayment 更新Payment信息
	UpdatePayment(*model.Payment) error
	// FindAll 查找所有Payment信息
	FindAll() ([]model.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

type PaymentRepository struct {
	db *gorm.DB
}

func (rep PaymentRepository) InitTable() error {
	return rep.db.CreateTable(&model.Payment{}).Error
}

func (rep PaymentRepository) FindPaymentByID(paymentID int64) (payment *model.Payment, err error) {
	payment = &model.Payment{}
	return payment, rep.db.Where("id = ?", paymentID).First(&payment).Error
}

func (rep PaymentRepository) CreatePayment(payment *model.Payment) (int64, error) {
	return payment.ID, rep.db.Create(payment).Error
}

func (rep PaymentRepository) DeletePaymentByID(paymentID int64) error {
	return rep.db.Where("id = ?", paymentID).Delete(&model.Payment{}).Error
}

func (rep PaymentRepository) UpdatePayment(payment *model.Payment) error {
	return rep.db.Updates(payment).Error
}

func (rep PaymentRepository) FindAll() (paymentAll []model.Payment, err error) {
	return paymentAll, rep.db.Find(&paymentAll).Error
}
