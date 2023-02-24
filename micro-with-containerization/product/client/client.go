package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/product/proto/product"
	"github.com/go-micro/plugins/v4/registry/consul"
	ow "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"io"
)

func main() {
	// 注册中心
	var consulRegistry = consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, ioCloser, err := common.NewTracer("go.micro.service.product.client", "127.0.0.1:6831")
	if err != nil {
		logrus.Fatal(err)
	}
	defer func(io io.Closer) {
		err = io.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(ioCloser)
	opentracing.SetGlobalTracer(t)

	// Create service
	srv := micro.NewService(
		micro.Name("go.micro.service.product.client"),
		micro.Version("1.0"),
		// 这里设置地址和余姚暴露的端口
		micro.Address("127.0.0.1:8082"),
		// 添加consul作为注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		micro.WrapClient(ow.NewClientWrapper(opentracing.GlobalTracer())),
	)

	productService := product.NewProductService("go.micro.service.product", srv.Client())

	productAdd := &product.ProductInfo{
		ProductName:        "go",
		ProductSku:         "book",
		ProductPrice:       100,
		ProductDescription: "go book",
		ProductCategoryId:  1,
		ProductImage: []*product.ProductImage{
			{
				ImageName: "go-book-image",
				ImageCode: "go-book-code",
				ImageUrl:  "http://www.google.com",
			}, {
				ImageName: "go-book-image1",
				ImageCode: "go-book-code1",
				ImageUrl:  "http://www.google.com1",
			}, {
				ImageName: "go-book-image2",
				ImageCode: "go-book-code2",
				ImageUrl:  "http://www.google.com2",
			},
		},
		ProductSize: []*product.ProductSize{
			{
				SizeName: "size-name",
				SizeCode: "size-code",
			},
		},
		ProductSeo: &product.ProductSeo{
			SeoTitle:       "go book seo title",
			SeoKeywords:    "go-book-seo-keywords",
			SeoDescription: "go-book-seo-description",
			SeoCode:        "go-book-seo-code",
		},
	}

	addProduct, err := productService.AddProduct(context.TODO(), productAdd)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(addProduct)
}
