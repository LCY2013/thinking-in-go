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
		micro.Name("micro.containerization.client"),
	)
	// 初始化
	service.Init()

	say := helloworld.NewHelloWorldService("micro.containerization.service", service.Client())
	resp, err := say.SayHello(context.TODO(), &helloworld.SayRequest{
		Message: "Hello world",
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"client": service.Name(),
		}).Error(err)
	}

	logrus.Info(resp)
}
