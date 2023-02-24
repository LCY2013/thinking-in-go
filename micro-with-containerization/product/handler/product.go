package handler

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/domain/service"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/product/proto/product"
)

type Product struct {
	productService service.IProductDataService
}

// AddProduct 添加商品信息
func (p Product) AddProduct(ctx context.Context, request *pb.ProductInfo, response *pb.ResponseProduct) error {
	productAdd := &model.Product{}
	if err := common.SwapTo(request, productAdd); err != nil {
		return err
	}

	productID, err := p.productService.AddProduct(productAdd)
	if err != nil {
		return err
	}

	response.ProductId = productID

	return nil
}

// FindProductByID 根据ID查找商品信息
func (p Product) FindProductByID(ctx context.Context, request *pb.RequestID, response *pb.ProductInfo) error {
	product, err := p.productService.FindProductByID(request.ProductId)
	if err != nil {
		return err
	}

	if err = common.SwapTo(response, product); err != nil {
		return err
	}

	return nil
}

// UpdateProduct 更新商品信息
func (p Product) UpdateProduct(ctx context.Context, request *pb.ProductInfo, response *pb.Response) error {
	productUpdate := &model.Product{}
	if err := common.SwapTo(request, productUpdate); err != nil {
		return err
	}

	if err := p.productService.UpdateProduct(productUpdate); err != nil {
		return err
	}

	response.Msg = "更新成功"
	return nil
}

// DeleteProductByID 删除商品信息
func (p Product) DeleteProductByID(ctx context.Context, productID *pb.RequestID, response *pb.Response) error {
	if err := p.productService.DeleteProduct(productID.ProductId); err != nil {
		return err
	}

	response.Msg = "删除成功"

	return nil
}

// FindAllProduct 查询所有商品信息
func (p Product) FindAllProduct(ctx context.Context, all *pb.RequestAll, response *pb.AllProduct) error {
	productAll, err := p.productService.FindAllProduct()
	if err != nil {
		return err
	}

	for _, product := range productAll {
		productInfo := &pb.ProductInfo{}
		err = common.SwapTo(product, productInfo)
		if err != nil {
			return err
		}

		response.ProductInfo = append(response.ProductInfo, productInfo)
	}

	return nil
}

// New Return a new handler
func New(productService service.IProductDataService) *Product {
	return &Product{
		productService: productService,
	}
}
