package table

import (
	"net/http"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

// TestDownload 确认http包的Get函数可以下载内容
// go test -v
func TestDownload(t *testing.T) {
	var urls = []struct {
		url        string
		statusCode int
	}{
		{
			"http://www.baidu.com/",
			http.StatusOK,
		}, {
			"http://tengine.taobao.org/",
			http.StatusOK,
		},
	}

	t.Log("Given the need to test downloading different content.")
	{
		for _, u := range urls {
			t.Logf("\t when checking \"%s\" for status code \"%d\"",
				u.url, u.statusCode)
			{
				resp, err := http.Get(u.url)
				if err != nil {
					t.Fatal("\t\tshould be able to Get the url.", ballotX, err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == u.statusCode {
					t.Logf("should have a \"%d\" status. %v", u.statusCode, checkMark)
				} else {
					t.Errorf("should have a \"%d\" status %v %v", u.statusCode, ballotX, resp.StatusCode)
				}
			}
		}
	}
}
