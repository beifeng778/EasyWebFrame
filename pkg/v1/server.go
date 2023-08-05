package v1

import (
	"net/http"
)

type Server interface {
	RouterTab
	Start(addr string) error
}

type sdkHTTPServer struct {
	Name    string
	handler Handler
	root    Filter
}

func (s *sdkHTTPServer) Route(method, path string, f handlerFuc) {
	s.handler.Route(method, path, f)
}

func (s *sdkHTTPServer) Start(addr string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		c := NewContext(writer, request)
		s.root(c)
	})
	return http.ListenAndServe(addr, nil)
}

func NewHTTPServer(name string, builder ...FilterBuilder) Server {
	handler := NewHandlerTree()
	var root = handler.ServeHTTP
	for i := len(builder) - 1; i >= 0; i-- {
		b := builder[i]
		root = b(root)
	}
	return &sdkHTTPServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}
