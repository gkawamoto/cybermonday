package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gkawamoto/cybermonday/config"
	"github.com/gkawamoto/cybermonday/handler"

	"gocloud.dev/server/requestlog"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	staticHandler := http.FileServer(http.Dir(conf.StaticDir))

	log.Printf("Serving static files from %s", conf.StaticDir)

	markdownHandler, err := handler.New(staticHandler, conf)
	if err != nil {
		log.Panic(err)
	}

	logHandler := requestlog.NewHandler(&logger{}, markdownHandler)

	s := &http.Server{
		Addr:    conf.Addr,
		Handler: logHandler,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic(err)
	}
}

type logger struct{}

func (l *logger) Log(e *requestlog.Entry) {
	log.Println(
		e.Request.Proto,
		e.Request.Method,
		e.Status,
		e.Request.URL.Path,
		e.ResponseBodySize,
	)
}
