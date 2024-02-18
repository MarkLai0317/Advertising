package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type ChiRouter struct {
	Router *chi.Mux
}

func NewChiAdapter() *ChiRouter {
	return &ChiRouter{Router: chi.NewRouter()}
}

func (chiRouter *ChiRouter) Get(path string, handlerFunc http.HandlerFunc) {
	chiRouter.Router.Get(path, handlerFunc)
}

func (chiRouter *ChiRouter) Post(path string, handlerFunc http.HandlerFunc) {
	chiRouter.Router.Post(path, handlerFunc)
}

// Use now accepts a pathPrefix. Middleware is applied only to paths with pathPrefix
func (chiRouter *ChiRouter) Use(pathPrefix string, middleware MiddlewareFunc) {
	if pathPrefix == "" { // If no pathPrefix is specified, apply middleware globally.
		chiRouter.Router.Use(func(next http.Handler) http.Handler {
			return middleware(next.ServeHTTP)
		})
	} else { // Apply middleware only to paths with the specified prefix.
		chiRouter.Router.Route(pathPrefix, func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return middleware(next.ServeHTTP)
			})
		})
	}
}

// default timeout for read and write
const TIMEOUT = 30 * time.Second

func (chiRouter *ChiRouter) ListenAndServe(address string, options ...ServerOption) error {
	server := &http.Server{
		Addr:    address,
		Handler: chiRouter.Router, // chi.Mux implements http.Handler
	}

	// Apply each provided ServerOption to the http.Server instance
	for _, option := range options {
		option(server)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Println("Stopping server")
		err := server.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	log.Println(fmt.Sprintf("Service listening on %s", address))

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
	// Start listening and serving HTTP requests

}

// WithReadTimeout configure http.Server parameter ReadTimeout
func WithReadTimeout(t time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.ReadTimeout = t
	}
}

// WithWriteTimeout configure http.Server parameter WriteTimeout
func WithWriteTimeout(t time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.WriteTimeout = t
	}
}
