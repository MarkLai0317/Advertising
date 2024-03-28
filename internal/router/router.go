package router

import "net/http"

// MiddlewareFunc definition

// ServerOption is a function that configures an http.Server
type ServerOption func(*http.Server)

// Update the WebFramework interface to include ListenAndServe
type WebFramework interface {
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Use(pathPrefix string, middleware func(http.Handler) http.Handler)
	ListenAndServe(address string, options ...ServerOption) error
	GetHandler() http.Handler
}
