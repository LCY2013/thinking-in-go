package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/logger"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/proto/gen/helloworld"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
)

func init() {
	logger.LoggerInit()
}

func main() {
	// 创建新服务
	service := micro.NewService(
		micro.Name("micro.containerization.service"),
	)

	// 初始化方法
	service.Init()
	// 注册服务
	helloworld.RegisterHelloWorldHandler(service.Server(), new(HelloWorld))
	// 运行服务
	if err := service.Run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"service": service.Name(),
		}).Error(err)
	}
}

type HelloWorld struct {
	// 需要实现方法
}

func (h HelloWorld) SayHello(ctx context.Context,
	request *helloworld.SayRequest,
	response *helloworld.SayResponse) error {
	response.Answer = "say: " + request.Message
	return nil
}
