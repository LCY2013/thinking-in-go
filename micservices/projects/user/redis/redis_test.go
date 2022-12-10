package redis

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

func TestGetRedisConn(t *testing.T) {
	tests := []struct {
		name    string
		want    redis.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRedisConn()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRedisConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRedisConn() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRedisLock(t *testing.T) {
	type args struct {
		key        string
		expireTime time.Duration
	}
	tests := []struct {
		name string
		args args
		want *redsync.Mutex
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRedisLock(tt.args.key, tt.args.expireTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRedisLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInitRedis go test -v -bench="TestInitRedis" -benchtime="60s"
func TestInitRedis(t *testing.T) {
	type args struct {
		host     string
		port     string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "clientTest",
			args: args{
				host:     "10.2.1.47",
				port:     "11111",
				password: "CFscrm!20210908#",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitRedis(tt.args.host, tt.args.port, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("InitRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i := 0; i < 100; i++ {
				go func(num int) {
					for {
						lock := GetRedisLock(fmt.Sprintf("%s-%d", "testConnect", num), time.Second*3)
						err := lock.Lock()
						if err != nil {
							t.Errorf("GetRedisLock() error = %v, wantErr %v", err, tt.wantErr)
							return
						}
						t.Logf("GetRedisLock() info = %v", lock)
						time.Sleep(time.Second * 3)
						lock.Unlock()
					}
				}(i)
			}
		})
	}
}
