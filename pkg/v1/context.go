package v1

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	R *http.Request
	W http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		R: r,
		W: w,
	}
}

func (c *Context) ReadJSON(req any) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJSON(code int, resp any) error {
	c.W.WriteHeader(code)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.W.Write(respJSON)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) OkJSON(data any) error {
	return c.WriteJSON(http.StatusOK, data)
}

func (c *Context) SystemInternalErr(data any) error {
	return c.WriteJSON(http.StatusInternalServerError, data)
}

func (c *Context) BadRequest(data any) error {
	return c.WriteJSON(http.StatusNotFound, data)
}
