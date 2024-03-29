package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
func helloHandler(c *gee.Context) {
	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
}

func indexErrHandler(c *gee.Context) {
	names := []string{"AMOS", "Jack", "Mike"}
	c.String(http.StatusOK, names[100])
}
func v2Handler(c *gee.Context) {
	// expect /hello/amos
	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // global midlleware

	r.GET("/", helloHandler)

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", v2Handler)
	}

	// index out of range for testing Recovery()
	r.GET("/panic", indexErrHandler)
	r.Run(":9999")

}
