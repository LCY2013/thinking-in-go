package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart-api/handler"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/cart-api/proto/cart-api"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/cart/proto/cart"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/select/roundrobin"
	ow "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"io"
	"net"
	"net/http"

	"github.com/micro/micro/v3/service/logger"
)

const (
	QPS = 100
)

func main() {
	// 注册中心
	var consulRegistry = consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, ioCloser, err := common.NewTracer("go.micro.api.cart-api", "127.0.0.1:6831")
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

	// 熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	// 启动熔断器
	common.GO(func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9096"), hystrixStreamHandler)
		if err != nil {
			logrus.Fatal(err)
		}
	})

	// Create cartService
	srv := micro.NewService(
		micro.Name("go.micro.api.cart-api"),
		micro.Version("1.0"),
		// 这里设置地址和余姚暴露的端口
		micro.Address("0.0.0.0:8086"),
		// 添加consul作为注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		micro.WrapHandler(ow.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		// 添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialize cartService
	srv.Init()

	// 构建 远程服务
	cartService := cart.NewCartService("go.micro.service.cart", srv.Client())

	// Register handler
	err = pb.RegisterCartApiHandler(srv.Server(), handler.New(cartService))
	if err != nil {
		logger.Fatal(err)
	}

	// Run cartService
	if err = srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(
	ctx context.Context,
	req client.Request,
	rsp any,
	opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// run 正常业务
		logrus.WithContext(ctx).Info(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		// 异常处理
		logrus.WithContext(ctx).Error(err)
		return err
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{
			Client: c,
		}
	}
}
