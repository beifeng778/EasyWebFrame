package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type HandlerOnTree struct {
	root *node
}

type node struct {
	path     string
	children []*node
	handler  handlerFuc
}

func (h *HandlerOnTree) ServeHTTP(c *Context) {
	path := c.R.URL.Path
	handler, ok := h.findRouter(path)
	if !ok {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)

}

func (h *HandlerOnTree) findRouter(path string) (handlerFuc, bool) {
	path = strings.Trim(path, "/")
	pattern := strings.Split(path, "/")
	cur := h.root
	for _, p := range pattern {
		child, ok := h.findChildPath(cur, p)
		if !ok {
			return nil, false
		}
		cur = child

	}
	if cur.handler == nil {
		return nil, false
	}
	return cur.handler, true
}
func (h *HandlerOnTree) Route(method, path string, f handlerFuc) {
	err := validRouter(path)
	if err != nil {
		fmt.Println("不允许注册改路由: ", path)
		return
	}
	path = strings.Trim(path, "/")
	pattern := strings.Split(path, "/")
	cur := h.root
	for index, p := range pattern {
		if matchChild, ok := h.childPath(cur, p); ok {
			cur = matchChild
		} else {
			h.createChildPath(cur, pattern[index:], f)
			return
		}
	}
}

func validRouter(path string) error {
	index := strings.Index(path, "*")
	if index > 0 {
		if index != len(path)-1 {
			fmt.Errorf("error:%v", errors.New("不允许注册该路由"))
		}
		if path[index-1] != '/' {
			fmt.Errorf("error:%v", errors.New("不允许注册该路由"))
		}
	}
	return nil
}

// 用于创建路由
func (h *HandlerOnTree) childPath(root *node, path string) (*node, bool) {
	for _, child := range root.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *HandlerOnTree) findChildPath(root *node, path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range root.children {
		if child.path == path && child.path != "*" {
			return child, true
		}
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
}

func (h *HandlerOnTree) createChildPath(root *node, path []string, handler handlerFuc) {
	cur := root
	for _, p := range path {
		nn := newNode(p)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handler
}

func newNode(p string) *node {
	return &node{
		path:     p,
		children: make([]*node, 0, 2),
	}
}

func NewHandlerTree() *HandlerOnTree {
	return &HandlerOnTree{root: &node{}}
}
