package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"fetch-rewards.com/mx/receipt-processor/pkg/handlers/rest"
	"fetch-rewards.com/mx/receipt-processor/pkg/service"
)

func main() {
	// init service
	srv := service.New()

	handler := muxHandlers(srv)

	// run http service & wait
	runService(handler)
}

func muxHandlers(service rest.ReceiptProcessorHandler) *http.ServeMux {
	handler := rest.MakeHTTPHandlers(service)

	// healthcheck for server
	healthcheckHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	// init handler with api and healthcheck
	muxHandler := http.NewServeMux()
	muxHandler.Handle("/healthcheck", healthcheckHandler)
	muxHandler.Handle("/", handler)

	return muxHandler
}

func runService(handler *http.ServeMux) {
	// init server
	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", 8000),
	}

	errc := make(chan error, 2)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logrus.WithField("port", 8000).Info("listening")
		errc <- server.ListenAndServe()
	}()

	logrus.WithFields(logrus.Fields{
		"port": 8000,
	}).Info("receipt-processor-api initialized...")

	logrus.WithFields(logrus.Fields{
		"reason": <-errc,
	}).Info("terminated")
}
