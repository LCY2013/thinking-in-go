package common

import (
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4/config"
)

type MysqlConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     int64  `json:"port"`
}

// GetMysqlFromConsul 获取mysql配置
func GetMysqlFromConsul(config config.Config, path ...string) *MysqlConfig {
	mysqlConfig := &MysqlConfig{}
	err := config.Get(path...).Scan(mysqlConfig)
	if err != nil {
		logrus.
			WithField("keyword", "GetMysqlFromConsul: 获取数据库配置失败！").
			Error(err)
	}
	return mysqlConfig
}
