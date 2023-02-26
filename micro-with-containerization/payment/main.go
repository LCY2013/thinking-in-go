package main

import (
	"fmt"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/repository"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/domain/service"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/handler"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/payment/proto/payment"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	ratelimit "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	ow "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
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
		common.Error(err)
	}

	// 注册中心
	var consulRegistry = consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 链路追踪
	t, ioCloser, err := common.NewTracer("go.micro.service.payment", "127.0.0.1:6831")
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

	// 暴露监控地址
	common.PrometheusBoot(9098)

	// Create service
	srv := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("1.0"),
		// 这里设置地址和余姚暴露的端口
		micro.Address("127.0.0.1:8088"),
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
		common.Error(err)
		return
	}
	defer func(db *gorm.DB) {
		err = db.Close()
		if err != nil {
			common.Error(err)
		}
	}(db)
	// 禁止复表
	db.SingularTable(true)

	//repository.NewPaymentRepository(db).InitTable()

	// 初始化数据服务
	paymentDataService := service.NewPaymentDataService(repository.NewPaymentRepository(db))

	// Register handler
	err = pb.RegisterPaymentHandler(srv.Server(), handler.New(paymentDataService))
	if err != nil {
		common.Error(err)
		return
	}

	// Run service
	if err = srv.Run(); err != nil {
		common.Error(err)
	}
}
