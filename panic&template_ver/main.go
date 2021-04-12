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
	"fmt"
	"html/template"
	"net/http"
	"time"

	"gee"
)

type student struct {
	Name string
	Age int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d",year, month, day)
}

func main() {
	r := gee.New()
	r.USE(gee.Logger()) //全部中间件，所有路由适用
	r.USE(gee.Recovery())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate" : FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static") //将路由中的assets替换为静态文件服务器目录./static

	student1 := &student{Name: "young",Age : 20}
	student2 := &student{Name:"amber", Age:19}
	r.GET("/", func (c *gee.Context){
		panicTest := []int{1,2,3}
		fmt.Println(panicTest[4])
		c.HTML(http.StatusOK,"css.tmpl", nil)
	})

	r.GET("/students", func(context *gee.Context) {
		context.HTML(http.StatusOK, "arr.tmpl", gee.H{
			//该部分key要与arr.tmpl中的查询字段一致。
			"title" : "gee",
			"students" : [2]*student{student1, student2},//[2]*student是指有两个*student类型的变量
		})
	})

	r.GET("/date", func(context *gee.Context) {
		context.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title" : "gee",
			"now" : time.Date(2021, 3, 8,0,0,0,0,time.UTC),
		})
	})

	r.Run(":9999")
}
