package main

import "fmt"

func main() {
	r := NewRouter()
	r.Use(test1)
	r.Use(test2)
	r.Use(test3)
	r.Run("123")

}
func test1(in string) string {
	return in + ":test1"
}
func test2(in string) string {
	return in + ":test2"
}
func test3(in string) string {
	return in + ":test3"
}

type middleware func(string) string
type Router struct {
	middlewareChain []func(string) string
	mux             []string
}

func NewRouter() *Router {
	return &Router{}
}
func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}
func (r *Router) Run(in string) {
	midlen := len(r.middlewareChain)
	r.mux = make([]string, midlen)
	for i := midlen - 1; i >= 0; i-- {
		fmt.Println("i:", i)
		res := r.middlewareChain[i](in)
		fmt.Println("res:", res)
		r.mux[i] = res
		fmt.Println("r.mux[i]:", r.mux[i])
	}
}
