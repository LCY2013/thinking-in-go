package sql

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sqlTestSuite struct {
	suite.Suite

	// 配置字段
	driver string
	dsn    string

	// 初始化字段
	db *sql.DB
}

func (s *sqlTestSuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) TearDownTest() {
	_, err := s.db.Exec("DELETE FROM test_model;")
	if err != nil {
		s.T().Fatal(err)
	}
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &sqlTestSuite{
		driver: "sqlite3",
		dsn:    "file:test.db?cache=shared&mode=memory",
	})
}
