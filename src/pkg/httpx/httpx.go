package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func Init(addr string, handler http.Handler) func() {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	startServer(srv)
	return func() {
		gracefulShutdown(srv)
	}
}

func startServer(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("http server listen err: %+v", err)
		}
	}()
}

func gracefulShutdown(srv *http.Server) {
	log.Infoln("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %+v", err)
	}
	select {
	case <-ctx.Done():
		fmt.Println("http exiting")
	default:
		fmt.Println("http server stopped")
	}
}
