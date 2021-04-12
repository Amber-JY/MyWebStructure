package main

/*
(1) index
curl -i http://localhost:9999/
HTTP/1.1 200 OK
Date: Sun, 01 Sep 2019 08:12:23 GMT
Content-Length: 19
Content-Type: text/html; charset=utf-8
<h1>Index Page</h1>

(2) v1
$ curl -i http://localhost:9999/v1/
HTTP/1.1 200 OK
Date: Mon, 12 Aug 2019 18:11:07 GMT
Content-Length: 18
Content-Type: text/html; charset=utf-8
<h1>Hello Gee</h1>*/

import (
	"log"
	"net/http"
	"time"

	"gee"
)

func onlyforv2() gee.HandlerFunc{
	return func (c *gee.Context){
		t := time.Now()//获得当前时间

		c.Next()//继续处理流程

		log.Printf("only for v2, %v", time.Since(t))
	}
}


func main() {
	r := gee.New()
	r.USE(gee.Logger()) //全部中间件，所有路由适用
	r.GET("/", func (c *gee.Context){
		c.HTML(http.StatusOK, "<h1>hello, gee<h1>")
	})

	v2 := r.Group("/v2")
	{
		v2.USE(onlyforv2())
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

	}

	r.Run(":9999")
}
