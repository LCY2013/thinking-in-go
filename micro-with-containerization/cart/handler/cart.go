package handler

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/domain/service"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/proto/cart"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
)

type Cart struct {
	cartDataService service.ICartDataService
}

// New Return a new handler
func New(cartDataService service.ICartDataService) *Cart {
	return &Cart{
		cartDataService: cartDataService,
	}
}

func (c Cart) AddCart(ctx context.Context, request *pb.CartInfo, response *pb.ResponseAdd) (err error) {
	cart := &model.Cart{}
	err = common.SwapTo(request, cart)
	if err != nil {
		return err
	}
	response.CartId, err = c.cartDataService.AddCart(cart)
	return err
}

func (c Cart) CleanCart(ctx context.Context, request *pb.Clean, response *pb.Response) error {
	err := c.cartDataService.CleanCart(request.UserId)
	if err != nil {
		return err
	}

	response.Msg = "购物车清空成功"
	return nil
}

func (c Cart) Incr(ctx context.Context, request *pb.Item, response *pb.Response) error {
	err := c.cartDataService.IncrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}

	response.Msg = "购物车添加成功"
	return nil
}

func (c Cart) Decr(ctx context.Context, request *pb.Item, response *pb.Response) error {
	err := c.cartDataService.DecrNum(request.Id, request.ChangeNum)
	if err != nil {
		return err
	}

	response.Msg = "购物车减少成功"
	return nil
}

func (c Cart) DeleteItemByID(ctx context.Context, cartID *pb.CartID, response *pb.Response) error {
	err := c.cartDataService.DeleteCart(cartID.Id)
	if err != nil {
		return err
	}

	response.Msg = "购物车删除成功"
	return nil
}

func (c Cart) GetAll(ctx context.Context, request *pb.CartFindAll, response *pb.CartAll) error {
	cartAll, err := c.cartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}

	for _, cart := range cartAll {
		ci := &pb.CartInfo{}
		err = common.SwapTo(cart, ci)
		if err != nil {
			return err
		}
		response.CartInfo = append(response.CartInfo, ci)
	}

	return nil
}
