package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/crontab/container"
	"github.com/LCY2013/thinking-in-go/crontab/master/configs"
	service "github.com/LCY2013/thinking-in-go/crontab/master/internal/job"
	webcontainer "github.com/LCY2013/thinking-in-go/crontab/master/internal/web/container"
	"github.com/LCY2013/thinking-in-go/crontab/master/internal/web/controller"
	_gin "github.com/LCY2013/thinking-in-go/crontab/third_party/gin"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"time"
)

func main() {
	app := fx.New(
		// Provide all the constructors we need, which teaches Fx how we'd like to
		// construct the *log.Logger, http.Handler, and *http.ServeMux types.
		// Remember that constructors are called lazily, so this block doesn't do
		// much on its own.
		fx.Provide(
			controller.NewJobController,
			webcontainer.NewContainer,
		),
		// Since constructors are called lazily, we need some invocations to
		// kick-start our application. In this case, we'll use Register. Since it
		// depends on an http.Handler and *http.ServeMux, calling it requires Fx
		// to build those types using the constructors above. Since we call
		// NewMux, we also register Lifecycle hooks to start and stop an HTTP
		// server.
		fx.Invoke(service.InitMgr),
		fx.Invoke(startServerApp),

		// This is optional. With this, you can control where Fx logs
		// its events. In this case, we're using a NopLogger to keep
		// our test silent. Normally, you'll want to use an
		// fxevent.ZapLogger or an fxevent.ConsoleLogger.
		/*fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),*/
	)

	// In a typical application, we could just use app.Run() here. Since we
	// don't want this example to run forever, we'll use the more-explicit Start
	// and Stop.

	if err := app.Start(context.Background()); err != nil {
		log.WithFields(log.Fields{
			"fx.start": "err",
		}).Error(err)
		return
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.WithFields(log.Fields{
			"fx.stop": "err",
		}).Error(err)
	}
}

// startServerApp App start
func startServerApp(webContainer *webcontainer.WebContainer) {
	app := container.NewApp(configs.Conf().AppName,
		container.BuildMultipleGinServe(configs.Conf().Serves),
		container.WithShutdownCallbacks(ContainerCallback))

	router, ok := container.GinEngineByServeName("master")
	if ok {
		router.POST("/job/save", _gin.Wrapper(webContainer.JobController.CreateJob))
		router.POST("/job/del", _gin.Wrapper(webContainer.JobController.DelJob))
	}

	app.StartAndServe()
}

// ContainerCallback App container stop callback
func ContainerCallback(ctx context.Context) {

}
