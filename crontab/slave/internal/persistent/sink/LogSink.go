package sink

import (
	"context"
	entity "github.com/LCY2013/thinking-in-go/crontab/domain"
	"github.com/LCY2013/thinking-in-go/crontab/lib/async"
	"github.com/LCY2013/thinking-in-go/crontab/slave/configs"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// LogSink 日志接受器
type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *entity.JobLog
	autoCommitChan chan *entity.LogBatch
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
		SetConnectTimeout(time.Duration(configs.Conf().MongoDB.ConnectionTimeout)*time.Millisecond).
		ApplyURI(configs.Conf().MongoDB.Uri)); err != nil {
		return
	}

	// 选择db和collection
	// 构建全局的logSink单例
	GLogSink = &LogSink{
		client:         client,
		logCollection:  client.Database(configs.Conf().Log.DB).Collection(configs.Conf().Log.Collection),
		logChan:        make(chan *entity.JobLog, 1000),
		autoCommitChan: make(chan *entity.LogBatch),
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
		log                 *entity.JobLog
		logBatch            *entity.LogBatch
		commitTimmer        *time.Timer
		autoCommitTimestamp = 1000
		autoCommitBatch     *entity.LogBatch
	)

	for {
		select {
		case log = <-logSink.logChan:
			// 把这条log写入mongodb
			// logSink.logCollection.insert(log)
			if logBatch == nil {
				logBatch = &entity.LogBatch{}
				// 让这个批次在某一个时间自动提交
				if configs.Conf().Log.AutoCommitTimestamp > 0 {
					autoCommitTimestamp = configs.Conf().Log.AutoCommitTimestamp
				}
				commitTimmer = time.AfterFunc(
					time.Duration(autoCommitTimestamp)*time.Millisecond,
					func(logBatch *entity.LogBatch) func() {
						// 这里不要直接操作chan batch，因为这里是另一个协程，容易造成并发问题
						return func() {
							logSink.autoCommitChan <- logBatch
						}
					}(logBatch))
			}

			// 添加这次的日志到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果达到了最大批次，就进行插入操作
			if len(logBatch.Logs) >= configs.Conf().Log.InsertBatchSize {
				// 保存日志信息
				logSink.savaBatchLog(logBatch)
				// 清空logBatch
				logBatch = nil
				// 取消定时器
				commitTimmer.Stop()
			}
		case autoCommitBatch = <-logSink.autoCommitChan: // 过期的批次
			// 判断过期自动提交的批次与当前的批次是否是同一批次
			if autoCommitBatch != logBatch {
				// 如果不是同一批次说明自动过期的批次已经提交过了，就进行下一次处理
				continue
			}
			// 把该批次写入数据存储
			logSink.savaBatchLog(autoCommitBatch)
			// 清空logBatch
			logBatch = nil
		}
	}

}

// savaBatchLog 批量写入日志
func (logSink *LogSink) savaBatchLog(batch *entity.LogBatch) {
	_, err := logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"savaBatchLog": "err",
		}).Error("get plugin info called")
		return
	}
}

// Append 追加日志信息到对应的数据源
func (logSink *LogSink) Append(jobLog *entity.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
		// 如果日志较多，队列满了就丢弃日志信息，或者在这里打印
	}
}
