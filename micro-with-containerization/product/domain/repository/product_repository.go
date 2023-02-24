package repository

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/domain/model"
	"github.com/jinzhu/gorm"
)

type IProductRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindProductByID 根据ProductID查找Product信息
	FindProductByID(int64) (*model.Product, error)
	// CreateProduct 创建Product
	CreateProduct(*model.Product) (int64, error)
	// DeleteProductByID 根据ProductID删除Product
	DeleteProductByID(int64) error
	// UpdateProduct 更新Product信息
	UpdateProduct(*model.Product) error
	// FindAll 查找所有Product信息
	FindAll() ([]model.Product, error)
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db: db}
}

type ProductRepository struct {
	db *gorm.DB
}

func (rep ProductRepository) InitTable() error {
	return rep.db.CreateTable(&model.Product{}, &model.ProductImage{}, &model.ProductSeo{}, &model.ProductSize{}).Error
}

func (rep ProductRepository) FindProductByID(productID int64) (product *model.Product, err error) {
	product = &model.Product{}
	return product, rep.db.
		Preload("ProductImage").
		Preload("ProductSize").
		Preload("ProductSeo").
		First(product, productID).Error
}

func (rep ProductRepository) CreateProduct(product *model.Product) (int64, error) {
	return product.ID, rep.db.Create(product).Error
}

func (rep ProductRepository) DeleteProductByID(productID int64) error {
	// 开启事物处理
	tx := rep.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 删除
	if err := tx.Unscoped().Where("id = ?", productID).Delete(&model.Product{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("image_product_id = ?", productID).Delete(&model.ProductImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("image_product_id = ?", productID).Delete(&model.ProductSize{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("image_product_id = ?", productID).Delete(&model.ProductSeo{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (rep ProductRepository) UpdateProduct(product *model.Product) error {
	return rep.db.Updates(product).Error
}

func (rep ProductRepository) FindAll() (productAll []model.Product, err error) {
	return productAll, rep.db.
		Preload("ProductImage").
		Preload("ProductSize").
		Preload("ProductSeo").
		Find(&productAll).Error
}
