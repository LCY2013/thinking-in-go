package handler

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/service"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/order/proto/order"
)

type Order struct {
	orderService service.IOrderDataService
}

// New Return a new handler
func New(orderService service.IOrderDataService) *Order {
	return &Order{
		orderService: orderService,
	}
}

func (o Order) GetOrderByID(ctx context.Context, request *pb.OrderID, response *pb.OrderInfo) error {
	order, err := o.orderService.FindOrderByID(request.OrderId)
	if err != nil {
		return err
	}

	if err = common.SwapTo(order, response); err != nil {
		return err
	}

	return nil
}

func (o Order) GetAllOrder(ctx context.Context, request *pb.AllOrderRequest, response *pb.AllOrder) error {
	allOrder, err := o.orderService.FindAllOrder()
	if err != nil {
		return err
	}

	for _, order := range allOrder {
		op := &pb.OrderInfo{}
		if err = common.SwapTo(order, op); err != nil {
			return err
		}
		response.OrderInfo = append(response.OrderInfo, op)
	}

	return nil
}

func (o Order) CreateOrder(ctx context.Context, request *pb.OrderInfo, response *pb.OrderID) error {
	orderAdd := &model.Order{}
	if err := common.SwapTo(request, orderAdd); err != nil {
		return err
	}

	orderID, err := o.orderService.AddOrder(orderAdd)
	if err != nil {
		return err
	}

	response.OrderId = orderID

	return nil
}

func (o Order) DeleteOrderByID(ctx context.Context, request *pb.OrderID, response *pb.Response) error {
	err := o.orderService.DeleteOrder(request.OrderId)
	if err != nil {
		return err
	}

	response.Msg = "删除订单信息成功"

	return nil
}

func (o Order) UpdateOrderPayStatus(ctx context.Context, request *pb.PayStatus, response *pb.Response) error {
	err := o.orderService.UpdatePayStatus(request.OrderId, request.PayStatus)

	if err != nil {
		return err
	}

	response.Msg = "更新订单支付信息成功"

	return nil
}

func (o Order) UpdateOrderShipStatus(ctx context.Context, request *pb.ShipStatus, response *pb.Response) error {
	err := o.orderService.UpdateShipStatus(request.OrderId, request.ShipStatus)

	if err != nil {
		return err
	}

	response.Msg = "更新订单发货信息成功"

	return nil
}

func (o Order) UpdateOrder(ctx context.Context, request *pb.OrderInfo, response *pb.Response) error {
	orderUpdate := &model.Order{}

	if err := common.SwapTo(request, orderUpdate); err != nil {
		return err
	}

	if err := o.orderService.UpdateOrder(orderUpdate); err != nil {
		return err
	}

	response.Msg = "更新订单信息成功"

	return nil
}
