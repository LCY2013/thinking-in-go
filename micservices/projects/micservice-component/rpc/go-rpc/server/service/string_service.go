/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-24
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package service

import (
	"errors"
	"time"
)

var (
	StrMaxSize = 100
	ErrMaxSize = errors.New("parameter to long")
)

// 定义rpc的请求结构体
type StringRequest struct {
	SA string
	SB string
}

// 定义rpc接口
type Service interface {
	// Concat sa and sb
	Concat(req StringRequest, ret *string) error
}

type StringService struct {
}

// 实现Service接口
func (stringService *StringService) Concat(req StringRequest, ret *string) error {
	// len(sa + sb) > StrMaxSize to ErrMaxSize
	if len(req.SA)+len(req.SB) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	time.Sleep(time.Second * 2)
	*ret = req.SA + req.SB
	return nil
}
