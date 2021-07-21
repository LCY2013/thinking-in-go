package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

//func TestRoutes(t *testing.T) {
//	tests := []struct {
//		name string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//		})
//	}
//}

func init() {
	Routes()
}

// TestSendJson 测试 /sendjson 内部服务端点
func TestSendJson(t *testing.T) {
	t.Log("Given the need to test SendJson endpoint.")
	{
		req, err := http.NewRequest("GET", "/sendjson", nil)
		if err != nil {
			t.Fatal("\t\tshould be able to make the get call.", ballotX, err)
		}
		t.Log("\t\tshould be able to make the get call.", checkMark)

		recorder := httptest.NewRecorder()

		// 服务默认的多路选择器（mux），ServeHTTP 方法来模仿外部客户端对 /sendjson 服务端点的请求
		// ServeHTTP 方法调用完成，http.ResponseRecorder 就包含来 SendJson 函数处理的响应
		http.DefaultServeMux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("\t\tshould receive a \"%d\" status. %v %v",
				http.StatusOK, ballotX, recorder.Code)
		}
		t.Log("\t\tshould receive a 200", checkMark)

		u := struct {
			Name  string
			Email string
		}{}

		if err := json.NewDecoder(recorder.Body).Decode(&u); err != nil {
			t.Fatal("\tshould decoded the response.", ballotX)
		}
		t.Log("\tshould decoded the response.", checkMark)

		if u.Name == "fufeng" {
			t.Log("\tshould have a Name.", checkMark)
		} else {
			t.Error("\tshould have a Name.", ballotX, u.Name)
		}

		if u.Email == "fufeng@email.com" {
			t.Log("\tshould have an Email.", checkMark)
		} else {
			t.Error("\tshould have an Email.", ballotX, u.Email)
		}
	}
}
