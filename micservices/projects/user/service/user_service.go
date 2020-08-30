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
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"time"
	"user/dao"
	"user/redis"
)

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVo struct {
	Username string
	Password string
	Email    string
}

var (
	ErrUserExisted  = errors.New("user is existed")
	ErrUserPassword = errors.New("email and password not match")
	ErrRegistering  = errors.New("email is registering")
)

type UserService interface {
	// 登陆接口
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)
	// 注册接口
	Register(ctx context.Context, vo *RegisterUserVo) (*UserInfoDTO, error)
}

type UserServiceImpl struct {
	userDao dao.UserDao
}

func MakeUserServiceImpl(userDao dao.UserDao) UserService {
	return &UserServiceImpl{
		userDao: userDao,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {
	user, err := userService.userDao.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, err
		} else {
			return nil, ErrUserPassword
		}
	} else {
		log.Printf("err : %s", err)
	}
	return nil, err
}

func (userService *UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVo) (*UserInfoDTO, error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err : %s", err)
		return nil, err
	}
	defer lock.Unlock()
	existUser, err := userService.userDao.SelectByEmail(vo.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}
		err = userService.userDao.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err
}
