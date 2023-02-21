package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"user/domain/repository"
	userSrv "user/domain/service"
	"user/handler"
	"user/proto/user"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name("user.service"),
		micro.Version("1.0"),
	)

	// initialise flags
	srv.Init()

	// 创建数据库连接
	db, err := gorm.Open("mysql",
		"root:123456@/micro?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logrus.Error(err)
		return
	}

	defer func(db *gorm.DB) {
		err = db.Close()
		if err != nil {
			logrus.Warn(err)
		}
	}(db)

	db.SingularTable(true)

	/*rp := repository.NewUserRepository(db)
	rp.InitTable()*/

	// 创建 UserService
	userService := userSrv.NewUserDataService(repository.NewUserRepository(db))

	// Register handler
	err = user.RegisterUserHandler(srv.Server(), handler.New(userService))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"service-name": "user.service",
		}).Error(err)
	}

	// Run service
	if err = srv.Run(); err != nil {
		logrus.Error(err)
		return
	}
}
