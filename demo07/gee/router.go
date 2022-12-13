package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将pattern拆分并返回一个parts字符串数组
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") // 按照"/"拆分

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由信息，并将路由信息映射到对应的处理函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern) // 获取路径数组

	key := method + "-" + pattern
	_, ok := r.roots[method] // 检查这个方法是否有对应的根节点，没有就创建一个根节点
	if !ok {                 // 如果不存在
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) // 将路径插入到trie树中去
	r.handlers[key] = handler                 // 将路由信息映射到对应的处理函数
}

// 获取路由信息
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 获取路径数组
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // 从根节点往下查询，直至查找到节点

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' { // 模糊匹配
				params[part[1:]] = searchParts[index] // 查找下去
			}
			if part[0] == '*' && len(part) > 1 { // 直接匹配到*
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
