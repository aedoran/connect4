package fiber

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type Map map[string]interface{}

// Handler defines a request handler used by middleware and routes.
type Handler func(*Ctx) error

// Config allows customizing app behavior.
type Config struct {
	ErrorHandler func(*Ctx, error) error
}

// App is a minimal HTTP application.
type App struct {
	mux        *http.ServeMux
	config     Config
	server     *http.Server
	middleware []Handler
}

// New creates a new App with the given config.
func New(cfg Config) *App {
	return &App{mux: http.NewServeMux(), config: cfg}
}

// Use registers a middleware handler.
func (a *App) Use(h Handler) {
	a.middleware = append(a.middleware, h)
}

// Get registers a handler for the given path.
func (a *App) Get(path string, h Handler) {
	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		chain := append([]Handler{}, a.middleware...)
		chain = append(chain, h)
		ctx := &Ctx{Request: r, ResponseWriter: w, chain: chain}
		if err := ctx.Next(); err != nil && a.config.ErrorHandler != nil {
			a.config.ErrorHandler(ctx, err)
		}
	})
}

// Listen starts the HTTP server.
func (a *App) Listen(addr string) error {
	a.server = &http.Server{Addr: addr, Handler: a.mux}
	return a.server.ListenAndServe()
}

// ShutdownWithContext gracefully stops the server.
func (a *App) ShutdownWithContext(ctx context.Context) error {
	if a.server == nil {
		return nil
	}
	return a.server.Shutdown(ctx)
}

// Test executes the app for testing purposes.
func (a *App) Test(req *http.Request, _ int) (*http.Response, error) {
	rr := httptest.NewRecorder()
	a.mux.ServeHTTP(rr, req)
	return rr.Result(), nil
}

// ErrServerClosed mirrors http.ErrServerClosed for compatibility.
var ErrServerClosed = http.ErrServerClosed

const StatusInternalServerError = http.StatusInternalServerError

// Ctx represents the request context passed to handlers.
type Ctx struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	chain          []Handler
	index          int
	statusCode     int
}

// Next executes the next handler in the chain.
func (c *Ctx) Next() error {
	if c.index >= len(c.chain) {
		return nil
	}
	h := c.chain[c.index]
	c.index++
	return h(c)
}

// Status sets the HTTP status code.
func (c *Ctx) Status(code int) *Ctx {
	c.statusCode = code
	return c
}

// JSON writes a JSON response.
func (c *Ctx) JSON(v interface{}) error {
	if c.statusCode == 0 {
		c.statusCode = http.StatusOK
	}
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(c.statusCode)
	return json.NewEncoder(c.ResponseWriter).Encode(v)
}

// Method returns the HTTP method.
func (c *Ctx) Method() string { return c.Request.Method }

// Path returns the request path.
func (c *Ctx) Path() string { return c.Request.URL.Path }

// Response provides minimal access to response status code.
type Response struct{ status int }

// Response returns a Response with the current status code.
func (c *Ctx) Response() *Response { return &Response{status: c.statusCode} }

// StatusCode reports the status code.
func (r *Response) StatusCode() int { return r.status }
