package main

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello jpc\n")
	})
	// 越界触发panic测试
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"jpc"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
