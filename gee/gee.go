package gee

import (
	"log"
	"net/http"
	"strings"
)

//首先定义了类型HandlerFunc,这是提供给框架用户的,用来定义路由映射的处理方法
//Engine 包含一张路由map router,key由请求方法和静态路由地址构成,例如GET-/、GET-/hello、POST-/hello,这样针对相同的路由,如果请求方法不同,可以映射不同的处理方法(Handler),value 是用户映射的处理方法。
//当用户调用(*Engine).GET()方法时,会将路由和处理方法注册到映射表 router 中,(*Engine).Run()方法,是 ListenAndServe 的包装。
//ServeHTTP 方法,解析请求的路径,查找路由映射表,如果查到,就执行注册的处理方法。如果查不到,就返回 404 NOT FOUND 。

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 支持嵌套
	engine      *Engine       // 所有路由组共享一个engine实例
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有的分组
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...) // 添加中间件
}

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

// addRoute函数，调用了group.engine.router.addRoute来实现了路由的映射
// 由于Engine从某种意义上继承了RouterGroup的所有属性和方法，因为 (*Engine).engine 是指向自己的。
// 这样实现，我们既可以像原来一样添加路由，也可以通过分组添加路由。
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

// GET
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

// POST
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 当我们接收到一个具体请求时，
// 要判断该请求适用于哪些中间件，在这里我们简单通过 URL 的前缀来判断。
// 得到中间件列表后，赋值给 c.handlers。
// handle 函数中，将从路由匹配得到的 Handler 添加到 c.handlers列表中，执行c.Next()。

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	engine.router.handle(c)
	c.handlers = middlewares
	engine.router.handle(c)
}
