package router_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MarkLai0317/Advertising/internal/router"

	"github.com/stretchr/testify/suite"
)

type ChiRouterUnitTestSuite struct {
	suite.Suite
}

func TestChiRouterUnitTestSuite(t *testing.T) {
	suite.Run(t, &ChiRouterUnitTestSuite{})
}

func (uts *ChiRouterUnitTestSuite) TestGet() {
	chiRouter := router.NewChiAdapter()
	helloWorldHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	}
	chiRouter.Get("/hello", helloWorldHandler)

	// Create a test server using the httptest package
	testServer := httptest.NewServer(chiRouter.Router)
	defer testServer.Close()

	// Make a GET request to the "/hello" endpoint
	resp, err := http.Get(testServer.URL + "/hello")

	uts.Equal(nil, err, "err getting /hello")
	if err != nil {
		log.Fatalf("err Get /hello: %s", err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("err Get /hello: %s", err)
	}

	uts.Equal(http.StatusOK, resp.StatusCode, "status code error")
	// Assert the status code and response body
	uts.Equal("Hello, world!", string(body), "response body error")

}

func (uts *ChiRouterUnitTestSuite) TestPost() {
	// Create test router with Post handler
	chiRouter := router.NewChiAdapter()
	helloWorldHandler := func(w http.ResponseWriter, r *http.Request) {
		// Read and return body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
	chiRouter.Post("/hello", helloWorldHandler)

	// Create a test server using the httptest package
	testServer := httptest.NewServer(chiRouter.Router)
	defer testServer.Close()

	// Make a POST request to the "/hello" endpoint
	resp, err := http.Post(testServer.URL+"/hello", "application/json", strings.NewReader(`{"message": "Hello, world!"}`))
	uts.Equal(nil, err, "err Posting /hello")
	if err != nil {
		uts.Fail("Failed to post /hello", "err Posting /hello: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		uts.Fail("Failed to read response body", "err reading response from /hello: %v", err)
	}
	// Assert the status code and response body
	uts.Equal(http.StatusOK, resp.StatusCode, "status code error")
	uts.Equal(`{"message": "Hello, world!"}`, string(body), "response body error")

}

type ChiRouterIntegrationTestSuite struct {
	suite.Suite
}

func TestChiRouterIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &ChiRouterIntegrationTestSuite{})
}
func (its *ChiRouterIntegrationTestSuite) TestUseMultilayerURL() {
	// define middleware
	testMiddleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add a custom header to the response
			w.Header().Set("X-Custom-Header-1", "middleware-applied")
			next.ServeHTTP(w, r) // Call ServeHTTP on the next handler
		})
	}

	testMiddleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add a custom header to the response
			w.Header().Set("X-Custom-Header-2", "middleware-applied")
			next.ServeHTTP(w, r) // Call ServeHTTP on the next handler
		})
	}

	// init router
	router := router.NewChiAdapter()
	// Applying middleware to a specific path prefix

	router.Use("/protected1", testMiddleware1)
	router.Use("/protected1", testMiddleware2)
	router.Use("/protected-prefix/protected2", testMiddleware1)
	router.Use("/protected-layer1", testMiddleware1)
	router.Use("/protected-layer1/protected-layer2", testMiddleware2)

	// Define a handler to use with and without the middleware
	helloWorldHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Handler response")
	}

	router.Get("/", helloWorldHandler)
	router.Get("/unprotected/resource", helloWorldHandler)
	router.Get("/protected1/resource", helloWorldHandler)
	router.Get("/protected-prefix/protected2", helloWorldHandler)
	router.Get("/protected-layer1/protected-layer2/resource", helloWorldHandler)

	testServer := httptest.NewServer(router.Router)
	defer testServer.Close()

	//Test the unprotected route that should not have middleware applied
	resp, err := http.Get(testServer.URL + "/")
	if err != nil {
		its.Fail("Failed to read response body", "err reading response /: %v", err)
	}
	its.Equal("", resp.Header.Get("X-Custom-Header-1"))
	its.Equal("", resp.Header.Get("X-Custom-Header-2"))

	// Test the unprotected route that should not have middleware applied
	resp, err = http.Get(testServer.URL + "/unprotected/resource")
	if err != nil {
		its.Fail("Failed to read response body", "err reading response /unprotected/resource: %v", err)
	}
	its.Equal("", resp.Header.Get("X-Custom-Header-1"))
	its.Equal("", resp.Header.Get("X-Custom-Header-2"))

	// Test the protected route that should have middleware1,2 applied
	resp, err = http.Get(testServer.URL + "/protected1/resource")
	if err != nil {
		its.Fail("Failed to read response body", "err reading response /protected1/resource: %v", err)
	}
	its.Equal("middleware-applied", resp.Header.Get("X-Custom-Header-1"))
	its.Equal("middleware-applied", resp.Header.Get("X-Custom-Header-2"))

	// Test the protected route that should have middleware1 applied
	resp, err = http.Get(testServer.URL + "/protected-prefix/protected2")
	if err != nil {
		its.Fail("Failed to read response body", "err reading response /protected-prefix/protected2: %v", err)
	}
	its.Equal("middleware-applied", resp.Header.Get("X-Custom-Header-1"))
	its.Equal("", resp.Header.Get("X-Custom-Header-2"))

	// Test middleware in different layer in the same route
	resp, err = http.Get(testServer.URL + "/protected-layer1/protected-layer2/resource")
	if err != nil {
		its.Fail("Failed to read response body", "err reading response /protected-layer1/protected-layer2/resource: %v", err)
	}
	its.Equal("middleware-applied", resp.Header.Get("X-Custom-Header-1"))
	its.Equal("middleware-applied", resp.Header.Get("X-Custom-Header-2"))
}
