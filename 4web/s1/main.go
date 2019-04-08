package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

/*
对 Go 语言来说，隐藏在框架之下的通常是 net/http 和 html/template 这两个标准库
net/http 标准库可以分为客户端和服务器两个部分，库中的结构和函数有些只支持客户端和服务器这两者中的一个，而有些则同时支持客户端和服务器：
	• Client 、 Response 、 Header 、 Request 和 Cookie 对客户端进行支持；
	• Server 、 ServeMux 、 Handler/HandleFunc 、 ResponseWriter 、 Header 、 Request和 Cookie 则对服务器进行支持
其中header、request、cookie两者都支持
*/
// DefaultServeMux
func main() {
	// ---------------------1
	// http.ListenAndServe("", nil)

	// ---------------2
	// server := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: nil,
	// }
	// server.ListenAndServe()

	// ------------------3 为服务器编写任何处理器
	// handler := myHandler{}
	// server := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: &handler,
	// }
	// server.ListenAndServe()

	// -------------多处理器、
	// handler1 := myHandler1{}
	// handler2 := myHandler2{}
	// server := http.Server{
	// 	Addr: ":8080",
	// 	// Handler: &handler,
	// }
	// // http.Handle("/h1", &handler1)
	// // http.Handle("/h2", &handler2)
	// http.DefaultServeMux.Handle("/h1", &handler1)
	// http.DefaultServeMux.Handle("/h2", &handler2)

	// server.ListenAndServe()
	// ----------不使用默认的多路复用器
	// mux := http.NewServeMux()
	// handler1 := myHandler1{}
	// handler2 := myHandler2{}
	// mux.Handle("/h1", &handler1)
	// mux.Handle("/h2", &handler2)
	// server := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: mux,
	// }

	// server.ListenAndServe()

	// -----------------处理器函数 将普通函数转换成处理器函数即可，他的ServerHTTP是调用自己的方法
	// server := http.Server{
	// 	Addr: ":8080",
	// 	// Handler: mux,
	// }
	// http.HandleFunc("/hello1", hello1)
	// http.HandleFunc("/hello2", hello2)

	// server.ListenAndServe()

	// -----------------处理器函数 将普通函数转换成处理器函数即可，他的ServerHTTP是调用自己的方法
	// mux := http.NewServeMux()
	// server := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: mux,
	// }
	// mux.HandleFunc("/hello1", hello1)
	// mux.HandleFunc("/hello2", hello2)

	// server.ListenAndServe()
	// -----------------处理器函数 将普通函数转换成处理器函数即可，他的ServerHTTP是执行自己
	// mux := http.NewServeMux()
	// server := http.Server{
	// 	Addr:    ":8080",
	// 	Handler: mux,
	// }
	// mux.HandleFunc("/hello1", http.HandlerFunc(hello1))
	// mux.HandleFunc("/hello2", http.HandlerFunc(hello2))

	// server.ListenAndServe()

	// 诸如日志记录、安全检查和错误处理这样的操作通常被称为横切关注点
	// type HandlerFunc func(ResponseWriter, *Request)  类型别名
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("/hello1", http.HandlerFunc(log(hello1)))
	mux.HandleFunc("/hello2", http.HandlerFunc(hello2))
	// mux.Handle("/a", t2(timeMiddleware(http.HandlerFunc(hello2))))
	r := NewRouter()
	r.Add(t2)
	r.Add(timeMiddleware)
	h2 := r.Use(http.HandlerFunc(hello2))
	mux.Handle("/a", h2)
	server.ListenAndServe()

}

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []func(http.Handler) http.Handler
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Add(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}
func (r *Router) Use(h http.Handler) http.Handler {
	var res = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		res = r.middlewareChain[i](res)
	}
	return res
}

func t2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
		fmt.Println("Handler function called -:" + name)
		next.ServeHTTP(w, r) //执行调用
	})
}
func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// http.Error(w, http.StatusText(403), 403)
			tStart := time.Now()
			next.ServeHTTP(w, r)
			tEnd := time.Since(tStart)
			fmt.Println("timeMiddleware:", tEnd)
		})
}

// 中间件
func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called -:" + name)
		h(w, r) //执行调用
	}
}

func hello1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello  111 ")
}
func hello2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello  2222 ")
}

type myHandler1 struct{}

func (my *myHandler1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World  111111 ")
}

type myHandler2 struct{}

func (my *myHandler2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World 222222222222 ")
}

// 在 Go 语言中 ， 一个处理器就是一个拥有 ServeHTTP 方法的接口 两个参数：第一个参数是一个 ResponseWriter 接口，而第二个参数则是一个指向 Reques t 结构的指针 。
// ServeHTTP(http.ResponseWriter, *http.Request)

// ServeMux 结构实现了ServeHTTP方法 ，DefaultServeMux 多路复用器是 ServeMux 结构的一个实例
// 根据请求的 URL 将请求重定向到不同的处理器
