package mongodb_usage

import (
	"context"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type mongoDbTestSuite struct {
	suite.Suite

	client      *mongo.Client
	url         string
	dialTimeout time.Duration

	err error
}

func (m *mongoDbTestSuite) SetupSuite() {
	// mongo客户端建立连接
	m.client, m.err = mongo.Connect(context.TODO(), options.Client().SetConnectTimeout(m.dialTimeout).ApplyURI(m.url))
	if m.err != nil {
		m.T().Error(m.err)
		return
	}
}

func (m *mongoDbTestSuite) TearDownTest() {
	if m.err != nil {
		// 关闭mongo连接
		err := m.client.Disconnect(context.TODO())
		if err != nil {
			return
		}
	}
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &mongoDbTestSuite{
		url:         "mongodb://root:123456@127.0.0.1:27017",
		dialTimeout: 5 * time.Second,
	})
}
