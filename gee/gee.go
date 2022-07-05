package gee

import (
	"fmt"
	"net/http"
)

//首先定义了类型HandlerFunc,这是提供给框架用户的,用来定义路由映射的处理方法
//Engine 包含一张路由map router,key由请求方法和静态路由地址构成,例如GET-/、GET-/hello、POST-/hello,这样针对相同的路由,如果请求方法不同,可以映射不同的处理方法(Handler),value 是用户映射的处理方法。
//当用户调用(*Engine).GET()方法时,会将路由和处理方法注册到映射表 router 中,(*Engine).Run()方法,是 ListenAndServe 的包装。
//ServeHTTP 方法,解析请求的路径,查找路由映射表,如果查到,就执行注册的处理方法。如果查不到,就返回 404 NOT FOUND 。

type HandleFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandleFunc
}

//constructor
func New() *Engine {
	return &Engine{router: make(map[string]HandleFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandleFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//GET
func (engine *Engine) GET(patttern string, handler HandleFunc) {
	engine.addRoute("GET", patttern, handler)
}

//POST
func (engine *Engine) POST(patttern string, handler HandleFunc) {
	engine.addRoute("POST", patttern, handler)
}

//Run
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Printf("404 NOT FOUND: %s\n", req.URL)
	}
}
