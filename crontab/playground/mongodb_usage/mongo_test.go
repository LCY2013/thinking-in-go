package mongodb_usage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
