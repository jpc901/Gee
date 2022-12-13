package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中支持间件
	parent      *RouterGroup  // 嵌套
	engine      *Engine       // 所有组共享一个Engine实例
}
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有的组
}

// New 构造函数 gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Use 将中间件应用到某个 Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Group 定义一个新的 RouterGroup
// 存储所有组共享相同的Engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup) // 添加一个分组
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	// 通过group去调用engine实例再去添加路由
	group.engine.router.addRoute(method, pattern, handler)
}

// GET GET请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST POST请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//ServeHTTP 函数当接收到一个具体请求时，要判断该请求适用于哪些中间件
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 在这里我们简单通过 URL 的前缀来判断。
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares // 得到中间件列表后，赋值给 c.handlers。
	engine.router.handle(c)
}
