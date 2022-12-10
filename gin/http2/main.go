package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
</head>
<body>
  <h1 style="color:red;">Welcome, fufeng!</h1>
</body>
</html>
`))

// main curl -v -k https://127.0.0.1:8080/welcome
func main() {
	logger := log.New(os.Stderr, "", 0)
	logger.Println("[WARNING] DON'T USE THE EMBED CERTS FROM THIS EXAMPLE IN PRODUCTION ENVIRONMENT, GENERATE YOUR OWN!")

	r := gin.Default()
	r.SetHTMLTemplate(html)

	r.GET("/welcome", func(context *gin.Context) {
		context.HTML(http.StatusOK, "https", gin.H{
			"status": "success",
		})
	})

	// Listen and Server in https://127.0.0.1:8080
	r.RunTLS(":8080", "http2/testdata/server.pem", "http2/testdata/server.key")
}
