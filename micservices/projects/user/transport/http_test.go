/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 Ramostear.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-14
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package transport

import (
	_ "context"
	"encoding/json"
	"flag"
	"net/http"
	"net/url"
	_ "reflect"
	"strings"
	"testing"
	"time"
	_ "user/endpoint"
)

func TestRegister(t *testing.T) {

	if !flag.Parsed() {
		flag.Parse()
	}

	args := flag.Args()
	postUrl := "http://127.0.0.1:19527/register"
	if len(args) > 0 {
		postUrl = args[0]
	}

	body := map[string]string{
		"email":    "fufeng@magic.com",
		"password": "fufeng",
		"username": "123456",
	}
	result, err := httpPost(postUrl, body)

	if err != nil {
		t.Errorf("http post err %s", err)
		t.FailNow()
	}

	t.Logf("result is %v", result)

}

func TestLogin(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}

	args := flag.Args()
	postUrl := "http://127.0.0.1:19527/login"
	if len(args) > 0 {
		postUrl = args[0]
	}
	body := map[string]string{
		"email":    "aoho@mail.com",
		"password": "aoho1",
	}
	result, err := httpPost(postUrl, body)

	if err != nil {
		t.Errorf("http post err %s", err)
		t.FailNow()
	}
	t.Logf("result is %v", result)
}

func httpPost(postUrl string, body map[string]string) (interface{}, error) {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}

	dataUrlVal := url.Values{}

	for k, v := range body {
		dataUrlVal.Add(k, v)
	}

	req, err := http.NewRequest("POST", postUrl, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//提交请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//读取返回值
	decode := json.NewDecoder(resp.Body)
	var result interface{}

	err = decode.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil

}
