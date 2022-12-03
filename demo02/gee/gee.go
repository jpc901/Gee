package gee

import "net/http"

// HandlerFunc 定义gee使用的请求处理程序
type HandlerFunc func(*Context)

// Engine 实现了ServeHTTP接口
type Engine struct {
	router *router
}

// New 是gee.Engine的构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 将路由和处理方法注册到映射表 *router* 中
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 定义添加GET请求的方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 定义添加POST请求的方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 定义启动HTTP服务的方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
