package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment-api/handler"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/payment-api/proto/payment-api"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/proto/payment"
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
)

func main() {
	// 注册中心
	var consulRegistry = consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, ioCloser, err := common.NewTracer("go.micro.api.payment-api", "127.0.0.1:6831")
	if err != nil {
		common.Error(err)
	}
	defer func(io io.Closer) {
		err = io.Close()
		if err != nil {
			common.Error(err)
		}
	}(ioCloser)
	opentracing.SetGlobalTracer(t)

	// 熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	// 启动熔断器
	common.GO(func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9089"), hystrixStreamHandler)
		if err != nil {
			common.Error(err)
		}
	})

	// 监控
	common.PrometheusBoot(9189)

	// Create cartService
	srv := micro.NewService(
		micro.Name("go.micro.api.payment-api"),
		micro.Version("1.0"),
		// 这里设置地址和需要暴露的端口
		micro.Address("0.0.0.0:8089"),
		// 添加consul作为注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		micro.WrapHandler(ow.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 作为服务端访问时生效
		micro.WrapClient(ow.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		// 添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialize cartService
	srv.Init()

	// 构建 远程服务
	paymentService := payment.NewPaymentService("go.micro.service.payment", srv.Client())

	// Register handler
	err = pb.RegisterPaymentApiHandler(srv.Server(), handler.New(paymentService))
	if err != nil {
		common.Error(err)
		return
	}

	// Run service
	if err = srv.Run(); err != nil {
		common.Error(err)
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
		common.Error(err)
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
