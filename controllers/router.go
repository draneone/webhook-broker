package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/imyousuf/webhook-broker/config"
	"github.com/imyousuf/webhook-broker/storage"
	"github.com/julienschmidt/httprouter"
)

var (
	apiRouter         *httprouter.Router
	listener          ServerLifecycleListener
	routerInitializer sync.Once
	server            *http.Server
	dataAccessor      storage.DataAccessor
)

// ServerLifecycleListener listens to key server lifecycle error
type ServerLifecycleListener interface {
	StartingServer()
	ServerStartFailed(err error)
	ServerShutdownCompleted()
}

// RequestLogger is a simple io.Writer that allows requests to be logged
type RequestLogger struct {
}

func (rLogger RequestLogger) Write(p []byte) (n int, err error) {
	log.Println(string(p))
	return len(p), nil
}

// ConfigureAPI configures API Server with interrupt handling
func ConfigureAPI(httpConfig config.HTTPConfig, listener ServerLifecycleListener, accessor storage.DataAccessor) *http.Server {
	dataAccessor = accessor
	routerInitializer.Do(func() {
		apiRouter = httprouter.New()
		setupAPIRoutes(apiRouter)
	})
	server = &http.Server{
		Handler:      handlers.LoggingHandler(RequestLogger{}, apiRouter),
		Addr:         httpConfig.GetHTTPListeningAddr(),
		ReadTimeout:  httpConfig.GetHTTPReadTimeout(),
		WriteTimeout: httpConfig.GetHTTPWriteTimeout(),
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		log.Println("Listening to http at -", httpConfig.GetHTTPListeningAddr())
		listener.StartingServer()
		if serverListenErr := server.ListenAndServe(); serverListenErr != nil {
			listener.ServerStartFailed(serverListenErr)
			log.Fatal(serverListenErr)
		}
	}()
	go func() {
		<-stop
		handleExit()
	}()
	return server
}

func handleExit() {
	log.Println("Shutting down the server...")
	serverShutdownContext, shutdownTimeoutCancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownTimeoutCancelFunc()
	server.Shutdown(serverShutdownContext)
	log.Println("Server gracefully stopped!")
	listener.ServerShutdownCompleted()
}

func setupAPIRoutes(apiRouter *httprouter.Router) {
	apiRouter.GET("/_status", Status)
}
