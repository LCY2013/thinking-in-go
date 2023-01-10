package sink

import (
	"context"
	entity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/slave/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// LogSink 日志接受器
type LogSink struct {
	client        *mongo.Client
	logCollection *mongo.Collection
	logChan       chan *entity.JobLog
}

var (
	// GLogSink 全局单例
	GLogSink *LogSink
)

// InitLogSink 初始化日志接收器
func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	// 建立mongodb连接
	if client, err = mongo.Connect(context.TODO(), options.Client().
		SetConnectTimeout(time.Duration(configs.Conf().MongoDB.ConnectionTimeout)).
		ApplyURI(configs.Conf().MongoDB.Uri)); err != nil {
		return
	}

	// 选择db和collection
	// 构建全局的logSink单例
	GLogSink = &LogSink{
		client:        client,
		logCollection: client.Database(configs.Conf().Log.DB).Collection(configs.Conf().Log.Collection),
		logChan:       make(chan *entity.JobLog, 1000),
	}

	// 启动mongodb处理协程
	async.GO(func() {
		GLogSink.writeLoop()
	})

	return
}

// writeLoop 日志存储协程
func (logSink *LogSink) writeLoop() {
	var (
		log *entity.JobLog
	)

	for {
		select {
		case log = <-logSink.logChan:
			// 把这条log写入mongodb
			// logSink.logCollection.insert(log)
		}
	}

}
