package handler

import (
	"context"
	"errors"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/payment-api/proto/payment-api"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/proto/payment"
	"github.com/plutov/paypal/v4"
	"strconv"
)

type PaymentApi struct {
	paymentService payment.PaymentService
}

// New Return a new handler
func New(paymentService payment.PaymentService) *PaymentApi {
	return &PaymentApi{
		paymentService: paymentService,
	}
}

var (
	ClientID = "AV5U9j3o1VSL4ihQH5kNr6icR3KUCYVfCpwNxuNl_cOssNyR6jMTAYmGJ9pO3DrZAr1ARvzx4pFkiI3G"
)

// PayPalRefund go.micro.api.payment-api 通过API向外暴露为/payment-api/paymentApi/payPalRefund, 接收http请求
// 即：//payment-api/paymentApi/payPalRefund请求会调用go.micro.service.payment 服务的 PayPalRefund 方法
func (p PaymentApi) PayPalRefund(ctx context.Context, request *pb.Request, response *pb.Response) error {
	// 验证 payment 支付通道是否赋值
	if err := isOK("payment_id", request); err != nil {
		response.StatusCode = 500
		return err
	}

	// 验证 退款号
	if err := isOK("refund_id", request); err != nil {
		response.StatusCode = 500
		return err
	}

	// 验证 退款金额
	if err := isOK("money", request); err != nil {
		response.StatusCode = 500
		return err
	}

	// 获取 支付通道ID
	pamentID, err := strconv.ParseInt(request.Get["payment_id"].Values[0], 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取支付通道信息
	paymentInfo, err := p.paymentService.FindPaymentByID(ctx, &payment.PaymentID{
		PaymentId: pamentID,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	// SID 获取, paymentInfo.PaymentSid
	// 支付模式
	status := paypal.APIBaseSandBox
	if paymentInfo.PaymentStatus {
		status = paypal.APIBaseLive
	}

	// 退款例子
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: request.Get["refund_id"].Values[0] + " micro 提醒你收款!",
			EmailMessage: request.Get["refund_id"].Values[0] + " 您有一个收款信息!",
			// 每笔转账都要唯一
			SenderBatchID: request.Get["refund_id"].Values[0],
		},
		Items: []paypal.PayoutItem{
			{
				RecipientType: "EMAIL",
				//RecipientWallet: "",
				Receiver: "sb-fhlb625157031@personal.example.com",
				Amount: &paypal.AmountPayout{
					// 币种
					Currency: "USD",
					Value:    request.Get["money"].Values[0],
				},
				Note:         request.Get["refund_id"].Values[0],
				SenderItemID: request.Get["refund_id"].Values[0],
			},
		},
	}

	// 创建支付客户端
	payPalClient, err := paypal.NewClient(ClientID, paymentInfo.PaymentSid, status)
	if err != nil {
		common.Error(err)
		return err
	}

	// 获取token
	_, err = payPalClient.GetAccessToken(ctx)
	if err != nil {
		common.Error(err)
		return err
	}

	paymentResult, err := payPalClient.CreateSinglePayout(ctx, payout)
	if err != nil {
		common.Error(err)
		return err
	}

	common.Info(paymentResult)

	response.Body = request.Get["refund_id"].Values[0] + "支付成功!"

	return nil
}

func isOK(key string, request *pb.Request) error {
	if _, ok := request.Get[key]; !ok {
		err := errors.New(key + "参数异常")
		common.Error(err)
		return err
	}
	return nil
}
