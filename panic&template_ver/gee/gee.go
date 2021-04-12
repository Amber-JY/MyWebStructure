package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		engine      *Engine       // all groups share a Engine instance
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // store all groups

		htmlTemplates *template.Template //加载模板到内存
		funcMap template.FuncMap //所有的自定义模板渲染函数
	}
)

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 定义分组调用中间件的函数
func (group *RouterGroup) USE(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares, middlewares...) //... => 切片展开
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	//接收到请求后判断适用于哪些中间件，并将适用的中间件保存在Context中的处理函数集handlers中
	for _, group := range engine.groups{
		if strings.HasPrefix(req.URL.Path, group.prefix){//使用前缀匹配查找适用的中间件
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

//
// @Description: 根据相对路径，映射到真实的文件，并将其返回，即可完成静态文件服务器
// @param relativePath 文件相对路径
// @param fs 依据根目录打开的文件系统
//
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc{
	absolutePath := path.Join(group.prefix, relativePath)//获得路由绝对路径

	fileserver := http.StripPrefix(absolutePath, http.FileServer(fs))//将前缀替换为文件目录主目录打开为文件服务
	//访问localhost:9999/assets/js/geektutu.js，
	//最终返回/usr/geektutu/blog/static/js/geektutu.js
	return func(c *Context) {
		file := c.Params["filepath"] //获得请求的文件名
		if _, err := fs.Open(file); err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		fileserver.ServeHTTP(c.Writer, c.Req)
	}
}

//
// @Description:静态文件系统接口，供分组路由调用
// @param relativePath 文件路由相对路径
// @param root	文件根目录
//
func (group *RouterGroup) Static(relativePath string, root string){
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")//获得通配路由

	group.GET(urlPattern, handler)
}


//提供设置自定义渲染函数和加载模板的方法
func (engine *Engine) SetFuncMap(funcmap template.FuncMap){
	engine.funcMap = funcmap
}

func (engine *Engine) LoadHTMLGlob(pattern string){
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}