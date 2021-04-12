package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//错误恢复中间件

func Recovery() HandlerFunc {
	return func(context *Context) {
		defer func() {
			if info := recover(); info != nil {
				message := fmt.Sprintf("%s", info)
				log.Printf("%s\n",trace(message))//日志输出调用栈信息
				context.JSON(http.StatusInternalServerError, "Internet Server Error.")
			}
		}()
		context.Next()
	}
}

func trace(message string) string{
	var pcs [32] uintptr//栈指针s
	n := runtime.Callers(3, pcs[:])//跳过前三个程序计数器，分别是Callers本身，上一层trace，再上一层的defer func
	var str strings.Builder//使用该方式建立字符串，提高性能，否则普通构建每次都需要拷贝
	str.WriteString("\nTraceback:")
	for _, pc := range pcs[:n]{
		fn := runtime.FuncForPC(pc)//获取对应的函数
		file, line := fn.FileLine(pc)//调用该函数的文件名和行号
		str.WriteString(fmt.Sprintf("\n\tfile:%s:%d", file, line))
	}
	return str.String()
}