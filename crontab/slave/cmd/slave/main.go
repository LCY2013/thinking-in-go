package main

import (
	"context"
	"github.com/LCY2013/thinking-in-go/crontab/container"
	"github.com/LCY2013/thinking-in-go/crontab/slave/configs"
	"github.com/LCY2013/thinking-in-go/crontab/slave/internal/scheduler"
	webcontainer "github.com/LCY2013/thinking-in-go/crontab/slave/internal/web/container"

	service "github.com/LCY2013/thinking-in-go/crontab/slave/internal/job"
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
			configs.Conf,
			webcontainer.NewContainer,
		),
		// Since constructors are called lazily, we need some invocations to
		// kick-start our application. In this case, we'll use Register. Since it
		// depends on an http.Handler and *http.ServeMux, calling it requires Fx
		// to build those types using the constructors above. Since we call
		// NewMux, we also register Lifecycle hooks to start and stop an HTTP
		// server.
		fx.Invoke(
			scheduler.InitScheduler,
			service.InitMgr,
			startServerApp,
		),
	)

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

	router, ok := container.GinEngineByServeName("salve")
	if ok {
		_ = router
	}

	app.StartAndServe()
}

// ContainerCallback App container stop callback
func ContainerCallback(ctx context.Context) {

}
