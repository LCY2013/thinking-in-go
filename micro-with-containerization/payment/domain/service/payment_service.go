package service

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/repository"
)

type IPaymentDataService interface {
	AddPayment(payment *model.Payment) (int64, error)
	DeletePayment(paymentID int64) error
	UpdatePayment(payment *model.Payment) error
	FindPaymentByID(paymentID int64) (payment *model.Payment, err error)
	FindAllPayment() ([]model.Payment, error)
}

// NewPaymentDataService 创建用户数据服务
func NewPaymentDataService(paymentRepository repository.IPaymentRepository) IPaymentDataService {
	return &PaymentService{
		paymentRepository: paymentRepository,
	}
}

type PaymentService struct {
	paymentRepository repository.IPaymentRepository
}

func (ser PaymentService) FindPaymentByID(paymentID int64) (payment *model.Payment, err error) {
	payment = &model.Payment{}
	return ser.paymentRepository.FindPaymentByID(paymentID)
}

func (ser PaymentService) AddPayment(payment *model.Payment) (int64, error) {
	// MQ
	return ser.paymentRepository.CreatePayment(payment)
}

func (ser PaymentService) DeletePayment(paymentID int64) error {
	return ser.paymentRepository.DeletePaymentByID(paymentID)
}

func (ser PaymentService) UpdatePayment(payment *model.Payment) error {
	// log
	return ser.paymentRepository.UpdatePayment(payment)
}

func (ser PaymentService) FindAllPayment() ([]model.Payment, error) {
	return ser.paymentRepository.FindAll()
}
