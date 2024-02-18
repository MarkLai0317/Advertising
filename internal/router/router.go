package router

import "net/http"

// MiddlewareFunc definition
type MiddlewareFunc func(handlerFunc http.HandlerFunc) http.HandlerFunc

// ServerOption is a function that configures an http.Server
type ServerOption func(*http.Server)

// Update the WebFramework interface to include ListenAndServe
type WebFramework interface {
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Use(pathPrefix string, middleware MiddlewareFunc)
	ListenAndServe(address string, options ...ServerOption) error
}
