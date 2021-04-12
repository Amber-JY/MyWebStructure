# Gee框架总结

## 文件架构

.├─gee
│      context.go
│      gee.go
│      go.mod
│      logger.go
│      recovery.go
│      router.go
│      trie.go
│
├─static
│  │  file1.txt
│  │
│  └─css
│          geektutu.css
│
└─templates
        arr.tmpl
        css.tmpl
        custom_func.tmpl

## Gee有哪些功能？

* 实现路由映射表，提供用户注册静态路由的方法

* 实现上下文，将路由独立，提供给用户大粒度的调用接口

* 使用Trie实现动态路由解析，支持***参数匹配*** 和***通配***，提供动态路由的注册方法

* 实现分组控制，细化路由粒度，增强框架可用性

* 支持中间件，可按分组增加拓展功能

* 实现静态文件资源服务，支持HTML Template渲染

* 实现错误处理机制，增强框架稳定性

  ***参数匹配*** ：例如 `/p/:lang/doc`，可以匹配 `/p/c/doc` 和 `/p/go/doc`

  ***通配*** ：例如`static/*filepath`可以匹配`static/js/jQuery.js`，也可以匹配`static/fav.ico`

## 框架调用流程

### 功能使用简述

* 使用New()生成框架实例
* 使用USE(MiddleFuncName())，添加中间件
* 使用SetFuncMap(), LoadHTMLGlob(), Static()设置模板函数表、HTML渲染模板库、静态文件服务器目录
* 使用GET()注册路由及处理函数，Group()生成分组对象

### 重要流程详解

#### 设计Context的原因

- 封装*http.Request和http.ResponseWriter， 简化接口的调用
- 支撑动态路由和中间件，保存相关参数
- 统一处理调用和处理函数的参数传递

#### 分组的意义

- 在应用场景下，需要对某一类相同路由进行相似的处理
  - 以`/admin`开头的路由需要鉴权
  - 以`/post`开头的路由匿名可访问
- 中间件给予框架的拓展能力，结合分组可以得到更好的应用，不同的分组可以应用不同的中间件，以实现分组的共同操作

## 一些问题

1. context.JSON(int, interface{})中，出现错误情形的不合理解决

   ```
   func (c *Context) JSON(code int, obj interface{}) {
   	c.SetHeader("Content-Type", "application/json")
   	c.Status(code)
   	//  c.StatusCode = code
   	//	c.Writer.WriteHeader(code)
   	encoder := json.NewEncoder(c.Writer)
   	if err := encoder.Encode(obj); err != nil {
   		http.Error(c.Writer, err.Error(), 500)
   		//  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//	w.Header().Set("X-Content-Type-Options", "nosniff")
   		//	w.WriteHeader(code)
   		//	fmt.Fprintln(w, error)
   	}
   }
   ```
   
   line3中已经设置了返回码，那么在`encoder.Encode(obj)`出现报错后，再修改`c.Writer`将无意义。较合理的处理方法应该是在错误处理中 报`panic`
   
2. trie中，前缀树路由冲突
   `c.Get("/a/:b", handler)`和`c.Get("/a/c",handler)`冲突，路由`/a/x`也将访问`/a/c`。此处也应当直接报`panic`