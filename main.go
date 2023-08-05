package main

import (
	v1 "EasyWebFrame/pkg/v1"
	"fmt"
	"net/http"
)

func main() {
	server := v1.NewHTTPServer("test01-server", v1.MetricFilter)
	server.Route("POST", "/singup", SingUp)
	server.Route("GET", "/hello/*", func(c *v1.Context) {
		c.W.WriteHeader(http.StatusOK)
		c.W.Write([]byte("hello!"))
	})
	server.Route("GET", "/hello/echo", func(c *v1.Context) {
		c.W.WriteHeader(http.StatusOK)
		c.W.Write([]byte("hello! good boy"))
	})
	server.Start(":8080")
}

func SingUp(ctx *v1.Context) {
	req := &signUpReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		fmt.Println(err)
		ctx.BadRequest(err)
		return
	}
	resp := &commonResponse{Data: 111}
	err = ctx.OkJSON(resp)
	if err != nil {
		fmt.Printf("write response error:%v", err)
		return
	}
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
