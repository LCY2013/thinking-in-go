package service

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/domain/repository"
)

type IProductDataService interface {
	AddProduct(product *model.Product) (int64, error)
	DeleteProduct(productID int64) error
	UpdateProduct(product *model.Product) error
	FindProductByID(productID int64) (product *model.Product, err error)
	FindAllProduct() ([]model.Product, error)
}

// NewProductDataService 创建用户数据服务
func NewProductDataService(productRepository repository.IProductRepository) IProductDataService {
	return &ProductService{
		productRepository: productRepository,
	}
}

type ProductService struct {
	productRepository repository.IProductRepository
}

func (ser ProductService) FindProductByID(productID int64) (product *model.Product, err error) {
	product = &model.Product{}
	return ser.productRepository.FindProductByID(productID)
}

func (ser ProductService) AddProduct(product *model.Product) (int64, error) {
	// MQ
	return ser.productRepository.CreateProduct(product)
}

func (ser ProductService) DeleteProduct(productID int64) error {
	return ser.productRepository.DeleteProductByID(productID)
}

func (ser ProductService) UpdateProduct(product *model.Product) error {
	// log
	return ser.productRepository.UpdateProduct(product)
}

func (ser ProductService) FindAllProduct() ([]model.Product, error) {
	return ser.productRepository.FindAll()
}
