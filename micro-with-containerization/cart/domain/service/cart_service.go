package service

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/domain/repository"
)

type ICartDataService interface {
	AddCart(cart *model.Cart) (int64, error)
	DeleteCart(cartID int64) error
	UpdateCart(cart *model.Cart) error
	FindCartByID(cartID int64) (cart *model.Cart, err error)
	FindAllCart(userID int64) ([]model.Cart, error)

	// CleanCart 情况购物车
	CleanCart(int64) error
	// IncrNum 增加商品数量
	IncrNum(int64, int64) error
	// DecrNum 减少商品数量
	DecrNum(int64, int64) error
}

// NewCartDataService 创建用户数据服务
func NewCartDataService(cartRepository repository.ICartRepository) ICartDataService {
	return &CartService{
		cartRepository: cartRepository,
	}
}

type CartService struct {
	cartRepository repository.ICartRepository
}

func (ser CartService) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	cart = &model.Cart{}
	return ser.cartRepository.FindCartByID(cartID)
}

func (ser CartService) AddCart(cart *model.Cart) (int64, error) {
	// MQ
	return ser.cartRepository.CreateCart(cart)
}

func (ser CartService) DeleteCart(cartID int64) error {
	return ser.cartRepository.DeleteCartByID(cartID)
}

func (ser CartService) UpdateCart(cart *model.Cart) error {
	// log
	return ser.cartRepository.UpdateCart(cart)
}

func (ser CartService) FindAllCart(userID int64) ([]model.Cart, error) {
	return ser.cartRepository.FindAll(userID)
}

func (ser CartService) CleanCart(userID int64) error {
	return ser.cartRepository.CleanCart(userID)
}

func (ser CartService) IncrNum(cartID int64, num int64) error {
	return ser.cartRepository.IncrNum(cartID, num)
}

func (ser CartService) DecrNum(cartID int64, num int64) error {
	return ser.cartRepository.DecrNum(cartID, num)
}
