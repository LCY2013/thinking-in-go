package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

// 编写基础示例

// go test -v -run="ExampleSendJson"
// SendJson 提供基础示例
// Example 作为函数名开头的后半部分存在的公开方法，这里会在 godoc 中生成对应的示例
// Test 作为函数名开头的后半部分公开的方法，这里表示的是测试示例函数
func ExampleSendJson() {
	r, _ := http.NewRequest("GET", "/sendjson", nil)
	recorder := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(recorder, r)

	var u struct {
		Name  string
		Email string
	}

	if err := json.NewDecoder(recorder.Body).Decode(&u); err != nil {
		log.Println("ERROR:", err)
	}

	// 使用 fmt 将结果写到 stdout 来检测输出
	fmt.Println(u)

}
