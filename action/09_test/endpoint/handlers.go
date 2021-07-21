package endpoint

import (
	"encoding/json"
	"net/http"
)

// Routes 为网络配置路由
func Routes() {
	http.HandleFunc("/sendjson", SendJson)
}

// SendJson 返回一个简单的 json 文件
func SendJson(writer http.ResponseWriter, request *http.Request) {
	u := struct {
		Name  string
		Email string
	}{
		"fufeng",
		"fufeng@email.com",
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(&u)
}
