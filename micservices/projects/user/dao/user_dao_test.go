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
package dao

import (
	"testing"
	"time"
)

func TestUserDaoImpl_Save(t *testing.T) {
	userDao := &UserDaoImpl{}

	err := InitMysql("127.0.0.1", "3306",
		"root", "123456", "go_project")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	user := &UserEntity{
		Username:  "fufeng",
		Password:  "123456",
		Email:     "magic@fufeng.com",
		CreatedAt: time.Now(),
	}
	err = userDao.Save(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new user ID : %d\n", user.ID)
}

func TestUserDaoImpl_SelectByEmail(t *testing.T) {
	userDao := &UserDaoImpl{}

	err := InitMysql("127.0.0.1", "3306",
		"root", "123456", "go_project")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	user, err := userDao.SelectByEmail("magic@fufeng.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result username is : %s\n", user.Username)
}
