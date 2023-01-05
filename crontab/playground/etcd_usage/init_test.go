package etcd_usage

import (
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

type etcdTestSuite struct {
	suite.Suite

	endpoints   []string
	dialTimeout time.Duration

	config clientv3.Config
	client *clientv3.Client
	err    error
}

func (s *etcdTestSuite) SetupSuite() {
	// 客户端配置
	s.config = clientv3.Config{
		Endpoints:   s.endpoints,
		DialTimeout: s.dialTimeout,
	}

	// 建立连接
	if s.client, s.err = clientv3.New(s.config); s.err != nil {
		s.T().Error(s.err)
		return
	}
}

func (s *etcdTestSuite) TearDownTest() {
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			s.T().Error(err)
			return
		}
	}
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &etcdTestSuite{
		endpoints:   []string{"127.0.0.1:2379"},
		dialTimeout: 5 * time.Second,
	})
}
