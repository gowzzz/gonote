package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func wshandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "aaa")
}
func main() {
	pool := x509.NewCertPool()
	caCertPath := "ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	fmt.Println("123")
	s := &http.Server{
		Addr: ":8080",
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	http.DefaultServeMux.HandleFunc("/ws", wshandler)
	if err := s.ListenAndServeTLS("./server.crt", "./server.key"); err != nil {
		fmt.Println("server err:", err)
	} else {
		fmt.Println("server ok")
	}
}
