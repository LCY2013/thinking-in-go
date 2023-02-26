package main

import (
	"fmt"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/repository"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/domain/service"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/order/handler"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/order/proto/order"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	ratelimit "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	ow "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/jinzhu/gorm"
	"github.com/micro/micro/v3/service/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	_ "gorm.io/driver/mysql"
	"io"
)

var (
	//qps = os.Getenv("QPS")

	QPS = 1000
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		logrus.Fatal(err)
	}

	// 注册中心
	var consulRegistry = consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, ioCloser, err := common.NewTracer("go.micro.service.order", "127.0.0.1:6831")
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

	// 暴露监控地址
	common.PrometheusBoot(9092)

	// Create service
	srv := micro.NewService(
		micro.Name("go.micro.service.order"),
		micro.Version("1.0"),
		// 这里设置地址和余姚暴露的端口
		micro.Address("127.0.0.1:8087"),
		// 添加consul作为注册中心
		micro.Registry(consulRegistry),
		// 绑定链路追踪
		micro.WrapHandler(ow.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		// 添加监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 获取mysql配置，路径中不带前缀
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	// 初始化mysql
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			mysqlInfo.User, mysqlInfo.Pwd, "127.0.0.1", "3306", mysqlInfo.Database))
	if err != nil {
		logrus.
			WithField("keyword", "mysql连接错误").
			Error(err)
		return
	}
	defer func(db *gorm.DB) {
		err = db.Close()
		if err != nil {
			logrus.
				WithField("keyword", "mysql关闭错误").
				Error(err)
		}
	}(db)
	// 禁止复表
	db.SingularTable(true)

	//repository.NewOrderRepository(db).InitTable()

	// 初始化数据服务
	orderDataService := service.NewOrderDataService(repository.NewOrderRepository(db))

	// Register handler
	err = pb.RegisterOrderHandler(srv.Server(), handler.New(orderDataService))
	if err != nil {
		return
	}

	// Run service
	if err = srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
