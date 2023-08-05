package v1

import (
	"net/http"
)

type RouterTab interface {
	Route(method, path string, f handlerFuc)
}
type Handler interface {
	ServeHTTP(c *Context)
	RouterTab
}
type handlerBaseOnMap struct {
	handlers map[string]handlerFuc
}

func (h *handlerBaseOnMap) ServeHTTP(c *Context) {
	key := h.key(c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("Not Found"))
	}
}
func (h *handlerBaseOnMap) Route(method, path string, f handlerFuc) {
	key := h.key(method, path)
	h.handlers[key] = f
}
func (h *handlerBaseOnMap) key(method, path string) string {
	return method + "#" + path
}

func NewHandlerBaseOnMap() Handler {
	return &handlerBaseOnMap{handlers: make(map[string]handlerFuc, 128)}
}
