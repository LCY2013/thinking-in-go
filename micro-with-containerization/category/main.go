package main

import (
	"fmt"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/repository"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/service"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/handler"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/category/proto/category"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	_ "gorm.io/driver/mysql"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		logrus.Error(err)
	}

	// 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// Create service
	srv := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("1.0"),
		// 这里设置地址和余姚暴露的端口
		micro.Address("127.0.0.1:8083"),
		// 添加consul作为注册中心
		micro.Registry(consulRegistry),
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

	repository.NewCategoryRepository(db).InitTable()

	// 初始化数据服务
	categoryDataService := service.NewCategoryDataService(repository.NewCategoryRepository(db))

	// Register handler
	pb.RegisterCategoryHandler(srv.Server(), handler.New(categoryDataService))

	// Run service
	if err = srv.Run(); err != nil {
		logrus.
			WithField("keyword", "服务启动失败").
			Fatal(err)
	}
}
