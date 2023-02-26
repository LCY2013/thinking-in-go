package handler

import (
	"context"
	"encoding/json"
	"errors"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/cart-api/proto/cart-api"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/proto/cart"
	"github.com/sirupsen/logrus"
	"strconv"
)

type CartApi struct {
	cartService cart.CartService
}

// New Return a new handler
func New(cartService cart.CartService) *CartApi {
	return &CartApi{
		cartService: cartService,
	}
}

// FindAll go.micro.api.cart-api 通过API向外暴露为/cart-api/carApi/findAll, 接收http请求
// 即：/cart-api/carApi/findAll请求会调用go.micro.service.cart 服务的FindAll方法
func (c CartApi) FindAll(ctx context.Context, request *pb.Request, response *pb.Response) error {
	logrus.WithContext(ctx).Info("cartapi/findAll 访问请求")
	if _, ok := request.Get["user_id"]; !ok {
		// response.StatusCode = 500
		return errors.New("参数异常")
	}

	userIdString := request.Get["user_id"].Values[0]
	logrus.WithContext(ctx).
		WithField("x-user-id", userIdString).
		Info("cartapi/findAll 访问请求")
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return err
	}

	// 获取购物车所有商品
	cartAll, err := c.cartService.GetAll(ctx, &cart.CartFindAll{
		UserId: userId,
	})
	if err != nil {
		return err
	}

	// 数据类型转化
	b, err := json.Marshal(cartAll)
	if err != nil {
		return err
	}

	response.StatusCode = 200
	response.Body = string(b)

	return nil
}
