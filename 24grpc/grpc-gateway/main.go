package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/tmc/grpc-websocket-proxy/wsproxy"

	"github.com/golang/glog"
	gw "github.com/gonote/24grpc/protoc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:50051", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	r := NewRouter()
	r.Add(timeMiddleware)

	return http.ListenAndServe(":8080", r.Use(mux))

	return http.ListenAndServe(":8080", r.Use(wsproxy.WebsocketProxy(mux)))
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		// next handler
		next.ServeHTTP(wr, r)
		timeElapsed := time.Since(timeStart)
		fmt.Println(timeElapsed)
	})
}

// ------------------------
type middleware func(http.Handler) http.Handler
type Router struct {
	middlewareChain []func(http.Handler) http.Handler
}

func NewRouter() *Router {
	return &Router{}
}
func (r *Router) Add(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}
func (r *Router) Use(h http.Handler) http.Handler {
	var mergedHandler = h
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}
	return mergedHandler
}
