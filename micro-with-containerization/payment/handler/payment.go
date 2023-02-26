package handler

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/service"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/proto/payment"
)

type Payment struct {
	paymentService service.IPaymentDataService
}

// New Return a new handler
func New(paymentService service.IPaymentDataService) *Payment {
	return &Payment{
		paymentService: paymentService,
	}
}

func (p Payment) AddPayment(ctx context.Context, request *pb.PaymentInfo, response *pb.PaymentID) error {
	payment := &model.Payment{}
	if err := common.SwapTo(request, payment); err != nil {
		common.Error(err)
		return err
	}

	paymentID, err := p.paymentService.AddPayment(payment)
	if err != nil {
		common.Error(err)
		return err
	}

	response.PaymentId = paymentID

	return nil
}

func (p Payment) UpdatePayment(ctx context.Context, request *pb.PaymentInfo, response *pb.Response) error {
	payment := &model.Payment{}
	if err := common.SwapTo(request, payment); err != nil {
		common.Error(err)
		return err
	}

	err := p.paymentService.UpdatePayment(payment)
	if err != nil {
		common.Error(err)
		return err
	}

	response.Msg = "更新账户信息成功"

	return nil
}

func (p Payment) DeletePayment(ctx context.Context, request *pb.PaymentID, response *pb.Response) error {
	return p.paymentService.DeletePayment(request.PaymentId)
}

func (p Payment) FindPaymentByID(ctx context.Context, request *pb.PaymentID, response *pb.PaymentInfo) error {
	payment, err := p.paymentService.FindPaymentByID(request.PaymentId)
	if err != nil {
		common.Error(err)
		return err
	}

	err = common.SwapTo(payment, response)
	if err != nil {
		common.Error(err)
		return err
	}

	return nil
}

func (p Payment) FindAllPament(ctx context.Context, arequestll *pb.All, response *pb.PaymentAll) error {
	payment, err := p.paymentService.FindAllPayment()
	if err != nil {
		common.Error(err)
		return err
	}

	for _, pt := range payment {
		paymentInfo := &pb.PaymentInfo{}
		err = common.SwapTo(pt, paymentInfo)
		if err != nil {
			common.Error(err)
			return err
		}
		response.PaymentInfo = append(response.PaymentInfo, paymentInfo)
	}

	return nil
}
