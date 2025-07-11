package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

// fiberResponseWriter implements http.ResponseWriter for Fiber
type fiberResponseWriter struct {
	c *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	return w.c.Write(b)
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.c.Status(statusCode)
}

// AdaptFiberToHTTP adapts a Fiber context to net/http.ResponseWriter and *http.Request
func AdaptFiberToHTTP(c *fiber.Ctx) (http.ResponseWriter, *http.Request) {
	w := &fiberResponseWriter{c: c}

	// Create a new http.Request from Fiber's fasthttp.Request
	httpReq := &http.Request{}

	// Method
	httpReq.Method = c.Method()

	// URL
	rawURL := c.OriginalURL()
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// Handle error, e.g., log it or return a default URL
		parsedURL = &url.URL{Path: "/"}
	}
	httpReq.URL = parsedURL

	// Header
	httpReq.Header = make(http.Header)
	c.Request().Header.VisitAll(func(key, value []byte) {
		httpReq.Header.Set(string(key), string(value))
	})

	// Body
	bodyBytes := c.Body()
	httpReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return w, httpReq
}


