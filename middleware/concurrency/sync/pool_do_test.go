package sync

import (
	"fmt"
	"sync"
	"testing"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *User) Reset() {
	if u == nil {
		return
	}
	u.Name = ""
	u.ID = 0
}

func (u *User) ChangeName(name string) {
	if u == nil {
		return
	}
	u.Name = name
}

func TestUserPool(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			return &User{}
		},
	}
	user1 := pool.Get().(*User)
	user1.ID = 1
	user1.Name = "user1"

	// 操作完放回去
	user1.Reset()
	pool.Put(user1)

	user2 := pool.Get().(*User)
	fmt.Println(user2)
}
