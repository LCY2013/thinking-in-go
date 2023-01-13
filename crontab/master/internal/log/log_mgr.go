package log

import (
	"context"
	entity "github.com/LCY2013/thinking-in-go/crontab/domain/job"
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Mgr 日志管理器
type Mgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	GLogMgr *Mgr
)

// InitLogMgr 初始化日志管理器
func InitLogMgr() (err error) {
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
	GLogMgr = &Mgr{
		client:        client,
		logCollection: client.Database(configs.Conf().Log.DB).Collection(configs.Conf().Log.Collection),
	}

	return
}

// ListLog 查看任务日志
func (mgr *Mgr) ListLog(name string, skip, limit int64) (logArr []*entity.JobLog, err error) {
	var (
		filter  *entity.JobLogFilter
		logSort *entity.SortLogByStartTime
		cursor  *mongo.Cursor
		jobLog  *entity.JobLog
	)

	// len(logArr
	logArr = make([]*entity.JobLog, 0)

	// 过滤条件
	filter = &entity.JobLogFilter{
		JobName: name,
	}

	// 排序条件，按任务开始时间倒排
	logSort = &entity.SortLogByStartTime{
		SortOrder: -1,
	}

	// 查询
	if cursor, err = mgr.logCollection.Find(context.TODO(), filter, options.Find().SetSort(logSort).SetSkip(skip).SetLimit(limit)); err != nil {
		return
	}

	// 延迟关闭游标
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			log.WithError(err).Error(err)
		}
	}(cursor, context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &entity.JobLog{}

		// 反序列化JSON
		if err = cursor.Decode(jobLog); err != nil {
			continue
		}

		logArr = append(logArr, jobLog)
	}

	return
}
