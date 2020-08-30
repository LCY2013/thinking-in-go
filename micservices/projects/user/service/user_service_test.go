/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-08-30
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package service

import (
	"context"
	"testing"
	"user/dao"
	"user/redis"
)

func TestUserServiceImpl_Login(t *testing.T) {
	err := dao.InitMysql("127.0.0.1", "3306",
		"root", "123456", "go_project")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDao: &dao.UserDaoImpl{},
	}

	user, err := userService.Login(context.Background(), "magic@fufeng.com", "123456")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)
}

func TestUserServiceImpl_Register(t *testing.T) {
	err := dao.InitMysql("127.0.0.1", "3306",
		"root", "123456", "go_project")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1", "6379", "123456")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDao: &dao.UserDaoImpl{},
	}

	user, err := userService.Register(context.Background(),
		&RegisterUserVo{
			Username: "lcy",
			Password: "123456",
			Email:    "lcy@fufeng.com",
		})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)
}
