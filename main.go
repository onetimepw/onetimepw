package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/onetimepw/onetimepw/build"
	"github.com/onetimepw/onetimepw/domain"
	"github.com/onetimepw/onetimepw/endpoint/web"
	"github.com/onetimepw/onetimepw/util/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type application struct {
	web *web.Endpoint
}

func newApplication(
	web *web.Endpoint,
) application {
	return application{
		web: web,
	}
}

func main() {
	var (
		err       error
		startTime time.Time
	)
	startTime = time.Now()
	build.StartTime = startTime

	var configPath = flag.String("config", "config.yml", "path to config file")

	flag.Parse()

	var config domain.Config
	err = configor.Load(&config, *configPath)
	if err != nil {
		fmt.Printf("can't load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(config.Env)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	app, err := InitApp(config)
	if err != nil {
		log.Fatal("Can't init app", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, eCtx := errgroup.WithContext(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	g.Go(func() error {
		zap.L().Info("Running fiber...")

		return app.web.Run(eCtx)
	})

	if err := g.Wait(); err != nil {
		log.Info("App terminated", zap.Error(err))
	}
}
