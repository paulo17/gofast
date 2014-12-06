// A fast web framework written in Go.
//
// Author: Vincent Composieux <vincent.composieux@gmail.com>

package gofast


import(
    "net/http"
    "log"
    "sort"
    "time"
)

type context struct {
    request    *request
    response   *response
    router     *router
    templating templating
}

// Creates a new context component instance
func NewContext() context {
    router     := NewRouter()
    templating := NewTemplating()

    return context{router: &router, templating: templating}
}

// Sets a HTTP request instance
func (c *context) SetRequest(req *http.Request, route route) {
    request := NewRequest(req, route)
    c.request = &request
}

// Returns a HTTP request component instance
func (c *context) GetRequest() *request {
    return c.request
}

// Sets a HTTP response instance
func (c *context) SetResponse(res http.ResponseWriter) {
    response := NewResponse(res)
    c.response = &response
}

// Returns a HTTP response component instance
func (c *context) GetResponse() *response {
    return c.response
}

// Returns a router component instance
func (c *context) GetRouter() *router {
    return c.router
}

// Returns a templating component instance
func (c *context) GetTemplating() templating {
    return c.templating
}

// Handles HTTP requests
func (c *context) Handle() {
    sort.Sort(RouteLen(c.GetRouter().GetRoutes()))

    http.ListenAndServe(PORT, c)
}

// Serves HTTP request by matching the correct route
func (c *context) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    for _, route := range c.GetRouter().GetRoutes() {
        if req.Method == route.method && route.pattern.MatchString(req.URL.Path) {
            c.SetRequest(req, route)

            res.Header().Set("Content-Type", "text/html; charset: utf-8")
            c.SetResponse(res)

            startTime := time.Now()
            route.handler(res, req)
            stopTime := time.Now()

            log.Printf("[%s] %q (time: %v)\n", req.Method, req.URL.String(), stopTime.Sub(startTime))
            break
        }
    }
}