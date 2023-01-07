package mongodb_usage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (m *mongoDbTestSuite) TestMongoDB() {
	var (
		database   *mongo.Database
		collection *mongo.Collection
	)

	// 选择数据库 my_db
	database = m.client.Database("dev")

	// 选择collection
	collection = database.Collection("dev_collection")

	_ = collection
}

// LogRecord 一条日志记录
type LogRecord struct {
	JobName   string    `json:"jobName" bson:"jobName"`     // 任务名称
	Command   string    `json:"command" bson:"command"`     // shell 命令
	Err       string    `json:"err" bson:"err"`             // 脚本错误信息
	Content   string    `json:"content" bson:"content"`     // 脚本输出
	TimePoint TimePoint `json:"timePoint" bson:"timePoint"` // 执行时间
}

// TimePoint 任务的执行时间点
type TimePoint struct {
	StartTime int64 `json:"startTime" bson:"startTime"`
	EndTime   int64 `json:"endTime" bson:"endTime"`
}

// TestMongoDBInsert mongodb insert operation
func (m *mongoDbTestSuite) TestMongoDBInsert() {
	var (
		database   *mongo.Database
		collection *mongo.Collection

		logRecord  *LogRecord
		insertResp *mongo.InsertOneResult
		docId      primitive.ObjectID

		err error
	)

	// 选择数据库 my_db
	database = m.client.Database("cron")

	// 选择collection
	collection = database.Collection("logs")

	// 插入记录(bson)
	logRecord = &LogRecord{
		JobName: "job0",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	if insertResp, err = collection.InsertOne(context.TODO(), logRecord); err != nil {
		m.T().Error(err)
		return
	}

	// _id: 默认生成一个全局唯一的ID，ObjectID: 12字节的二进制
	docId = insertResp.InsertedID.(primitive.ObjectID)
	m.T().Logf("自增ID: %s", docId.Hex())
}

// FindByJobName LogRecord filter by job name
type FindByJobName struct {
	JobName string `json:"jobName" bson:"jobName"` // JobName赋值为job0
}

// TestMongoDBFind mongo find records
func (m *mongoDbTestSuite) TestMongoDBFind() {
	var (
		database   *mongo.Database
		collection *mongo.Collection

		logRecord *LogRecord
		cond      *FindByJobName
		cursor    *mongo.Cursor

		err error
	)

	// 选择数据库 my_db
	database = m.client.Database("cron")

	// 选择collection
	collection = database.Collection("logs")

	// 按照jobName 字段过滤，找出jobName=job0的记录
	cond = &FindByJobName{
		JobName: "job0",
	}

	// 查询(过滤 + 翻页)
	if cursor, err = collection.Find(context.TODO(), cond, options.Find().SetSkip(0).SetLimit(2)); err != nil {
		m.T().Error(err)
		return
	}

	// 释放游标
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err = cursor.Close(ctx)
		if err != nil {
			m.T().Error(err)
			return
		}
	}(cursor, context.TODO())

	// 遍历结果集
	for cursor.Next(context.TODO()) {
		logRecord = &LogRecord{}

		// 反序列化
		if err = cursor.Decode(logRecord); err != nil {
			m.T().Error(err)
			return
		}

		// 把日志行打印出来
		m.T().Logf("record: %+v", logRecord)
	}
}

// TimeBeforeCond 时间条件
// {"$lt": 当前时间}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// DeleteCond {"timePoint.startTime": {"$lt": 当前时间}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

// TestMongoDBFind mongo delete records
func (m *mongoDbTestSuite) TestMongoDBDelete() {
	var (
		database   *mongo.Database
		collection *mongo.Collection
		delResp    *mongo.DeleteResult

		delCond *DeleteCond

		err error
	)

	// 选择数据库 my_db
	database = m.client.Database("cron")

	// 选择collection
	collection = database.Collection("logs")

	// 删除开始时间早于当前时间的所有日志($lt是less than)
	// delete({"timePoint.startTime": {"$lt": 当前时间}})
	delCond = &DeleteCond{
		BeforeCond: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}

	// 执行删除
	if delResp, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		m.T().Error(err)
		return
	}

	m.T().Logf("删除行数: %d", delResp.DeletedCount)
}
