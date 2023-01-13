package configs

import (
	"github.com/LCY2013/thinking-in-go/crontab/container"
	_etcd "github.com/LCY2013/thinking-in-go/crontab/third_party/etcd"
	_mongo "github.com/LCY2013/thinking-in-go/crontab/third_party/mongo"
	"github.com/LCY2013/thinking-in-go/crontab/tools"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	AppName string               `json:"appName"`
	Serves  []container.Serve    `json:"serves"`
	Etcd    _etcd.EtcdConfig     `json:"etcd"`
	MongoDB _mongo.MongoDBConfig `json:"mongodb"`
	Log     struct {
		DB         string `json:"db"`
		Collection string `json:"collection"`
	} `json:"log"`
}

var (
	once sync.Once
	conf *Config
)

// initConfigArg 初始化一些命令行参数
// serve -port 8080
func initConfigArg() {
	// 对数组没用
	pflag.Int("serves.readTimeOut", 5000, "server read time out")
	pflag.Int("serves.writeTimeOut", 5000, "server write time out")

	pflag.Parse()
}

// readInConfig 开始初始化整个配置
func readInConfig() error {
	var (
		pwd string
		err error
	)
	initConfigArg()
	if err = viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}
	viper.SetConfigName(*pflag.String("config.name", "config", "config name"))         // name of config file (without extension)
	viper.SetConfigType(*pflag.String("config.extension", "yaml", "config extension")) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/crontab/")                                               // path to look for the config file in
	viper.AddConfigPath("$HOME/.crontab")                                              // call multiple times to add many search paths
	viper.AddConfigPath(".")                                                           // optionally look for config in the working directory
	viper.AddConfigPath("./configs")                                                   // optionally look for config in the working directory
	viper.AddConfigPath("../configs")                                                  // optionally look for config in the working directory
	viper.AddConfigPath("../../configs")                                               // optionally look for config in the working directory
	viper.AddConfigPath("../../../configs")                                            // optionally look for config in the working directory
	if pwd, err = tools.Pwd(); err == nil {
		viper.AddConfigPath(pwd) // optionally look for config in the working directory
	}

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		return err
	}
	conf = &Config{}

	return viper.Unmarshal(conf, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
}

// Conf 直接获取conf
func Conf() *Config {
	if conf == nil {
		once.Do(func() {
			if err := readInConfig(); err != nil {
				log.WithFields(log.Fields{
					"initConfig": "Conf",
				}).Error(err)
			}
		})
	}
	return conf
}

// SafeConf 获取提示信息
func SafeConf() (*Config, error) {
	var err error
	if conf == nil {
		once.Do(func() {
			err = readInConfig()
		})
	}
	return conf, err
}
