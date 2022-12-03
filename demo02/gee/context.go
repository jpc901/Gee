package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 给map[string]interface{}起了一个别名gee.H 用来存放各种响应数据gin中源码也是类似
type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	// 响应信息
	StatusCode int
}

// 返回一个Context实例
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm 获取Query和PostForm参数的方法。
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取Query和PostForm参数的方法。
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置响应头状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头返回类型
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 构造String响应方法
func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value)))
}

// JSON 构造JSON响应方法
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 构造Data响应方法
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 构造HTML响应方法
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
