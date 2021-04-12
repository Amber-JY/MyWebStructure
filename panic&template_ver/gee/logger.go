package gee

import (
	"log"
	"time"
)

//定义日志输出中间件
//返回处理程序
func Logger() HandlerFunc {
	return func (c *Context){
		t := time.Now()//获得当前时间

		c.Next()//继续处理流程

		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))//日志输出，请求目的的返回状态和耗费时间


	}
}