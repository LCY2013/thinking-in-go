package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"time"
)

// HTML 渲染
func formatAsDate(t time.Time) string {
	date, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", date, month, day)
}

func main() {
	router := gin.Default()
	router.Delims("{[{", "}]}")
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

	router.LoadHTMLGlob("template/testdata/*")
	//router.LoadHTMLFiles("./testdata/raw.tmpl")

	router.GET("/raw", func(context *gin.Context) {
		context.HTML(http.StatusOK, "raw.tmpl", gin.H{
			"now": time.Date(2021, 07, 01, 0, 0, 0, 0, time.Local),
		})
	})

	router.Run(":8080")
}
