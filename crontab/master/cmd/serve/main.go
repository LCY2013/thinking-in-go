package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/crontab/container"
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	/*router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	err := router.Run(fmt.Sprintf(":%d", configs.Conf()))
	if err != nil {
		log.WithFields(log.Fields{
			"serve": "gin router",
		}).Error(err)
		return
	}*/

	/*log.WithFields(log.Fields{
		"initConfig": "Conf",
	}).Printf("%+v", *configs.Conf())*/

	app := container.NewApp(configs.Conf().AppName,
		container.BuildMultipleGinServe(configs.Conf().Serves),
		container.WithShutdownCallbacks(StoreCacheToDBCallback))

	router, ok := container.GinEngineByServeName("")
	if ok {
		router.GET("/index", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	app.StartAndServe()
}

func StoreCacheToDBCallback(ctx context.Context) {

}
