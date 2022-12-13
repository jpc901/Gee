package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果height已经到达字符串数组末尾直接返回
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 当前节点对应的part
	part := parts[height]
	child := n.matchChild(part) // 第一个匹配成功的节点
	if child == nil {
		// 新建一个结点加入到当前节点的孩子数组中去
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1) // 继续往下找
}

// 查找节点
func (n *node) search(parts []string, height int) *node {
	// 如果是height已经遍历完parts或者n.part包含字符*前缀就直接返回当前节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part) // 当前节点的孩子节点的数组

	for _, child := range children {
		result := child.search(parts, height+1) // 递归往下找节点
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
