package http_server

import (
	"net/http"
	"testing"
)

// TestDownload 确认 http 包的 Get 函数可以下载内容，并且内容可以被序列化和反序列化
// 在没有网络的情况下，通过 mock 函数执行模拟，让测试代码执行通过
func TestDownload(t *testing.T) {
	statusCode := http.StatusOK

	server := mockServer()
	defer server.Close()

	t.Log("Given the need to test downloading content.")
	{
		t.Logf("\twhen checking \"%s\" for status code \"%d\"",
			server.URL, statusCode)
		{
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatal("\t\tshould be able to make the Get call.",
					ballotX, err)
			}
			t.Log("\t\tshould be able to make the Get call.", checkMark)

			defer resp.Body.Close()

			if resp.StatusCode != statusCode {
				t.Fatalf("\t\tshould receive a \"%d\" status. %v %v",
					statusCode, ballotX, resp.StatusCode)
			}
			t.Logf("\t\tshould receive a \"%d\" status. %v", statusCode, checkMark)
		}
	}

}
