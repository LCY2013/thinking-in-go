package repository

import (
	"errors"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/domain/model"
	"github.com/jinzhu/gorm"
)

type ICartRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindCartByID 根据CartID查找Cart信息
	FindCartByID(int64) (*model.Cart, error)
	// CreateCart 创建Cart
	CreateCart(*model.Cart) (int64, error)
	// DeleteCartByID 根据CartID删除Cart
	DeleteCartByID(int64) error
	// UpdateCart 更新Cart信息
	UpdateCart(*model.Cart) error
	// FindAll 查找所有Cart信息
	FindAll(int64) ([]model.Cart, error)

	// CleanCart 情况购物车
	CleanCart(int64) error
	// IncrNum 增加商品数量
	IncrNum(int64, int64) error
	// DecrNum 减少商品数量
	DecrNum(int64, int64) error
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &CartRepository{db: db}
}

type CartRepository struct {
	db *gorm.DB
}

func (rep CartRepository) InitTable() error {
	return rep.db.CreateTable(&model.Cart{}).Error
}

func (rep CartRepository) FindCartByID(cartID int64) (cart *model.Cart, err error) {
	cart = &model.Cart{}
	return cart, rep.db.Where("id = ?", cartID).First(&cart).Error
}

func (rep CartRepository) CreateCart(cart *model.Cart) (int64, error) {
	db := rep.db.FirstOrCreate(cart, model.Cart{
		ProductID: cart.ProductID,
		SizeID:    cart.SizeID,
		UserID:    cart.UserID,
	})
	if db.Error != nil {
		return 0, db.Error
	}
	if db.RowsAffected == 0 {
		return 0, errors.New("购物车新增失败")
	}
	return cart.ID, nil
}

func (rep CartRepository) DeleteCartByID(cartID int64) error {
	return rep.db.Where("id = ?", cartID).Delete(&model.Cart{}).Error
}

func (rep CartRepository) UpdateCart(cart *model.Cart) error {
	return rep.db.Updates(cart).Error
}

func (rep CartRepository) FindAll(userID int64) (cartAll []model.Cart, err error) {
	return cartAll, rep.db.Where("user_id = ?", userID).Find(&cartAll).Error
}

func (rep CartRepository) CleanCart(userID int64) error {
	return rep.db.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}

func (rep CartRepository) IncrNum(cartID int64, num int64) error {
	cart := &model.Cart{}
	return rep.db.Model(cart).UpdateColumn("num", gorm.Expr("num + ?", num)).Error
}

func (rep CartRepository) DecrNum(cartID int64, num int64) error {
	cart := &model.Cart{
		ID: cartID,
	}
	db := rep.db.Model(cart).Where("num >= ?", num).
		UpdateColumn("num", gorm.Expr("num - ?", num))
	if db.Error != nil {
		return db.Error
	}

	if db.RowsAffected == 0 {
		return errors.New("减少购物车数量失败")
	}

	return nil
}
