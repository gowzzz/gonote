package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// proxyAddr := "http://123.132.232.254:61017/"
	// proxyAddr := "http://112.95.224.58:8118/"
	// 39.108.168.155	8118
	// 111.77.197.69	9999
	proxyAddr := "http://111.77.197.69:9999/"
	httpUrl := "https://m.gmw.cn/baijia/2019-04/02/1300272610.html"
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("proxy:", proxy)
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	fmt.Println(netTransport)
	httpClient := &http.Client{
		Timeout:   time.Second * 5,
		Transport: netTransport,
	}
	res, err := httpClient.Get(httpUrl)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Println(err)
		return
	}
	c, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(c))
}
