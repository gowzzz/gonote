package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

/*
Colly使用Golang的默认http客户端作为网络层。
可以通过更改默认的HTTP roundtripper来调整HTTP选项。


可配置的环境变量
ALLOWED_DOMAINS (可访问的域名，是一个切片)
CACHE_DIR (string)
DETECT_CHARSET (y/n)
DISABLE_COOKIES (y/n)
DISALLOWED_DOMAINS (禁止访问的域名，是一个切片)
IGNORE_ROBOTSTXT (y/n)
MAX_BODY_SIZE (int)
MAX_DEPTH (int - 0 意味着无限制)
PARSE_HTTP_ERROR_RESPONSE (y/n)
USER_AGENT (string)


OnRequest 在请求之前调用
OnError 在请求中出现错误时调用
OnResponse 响应接收到之后调用
OnHTML OnResponse 正确执行后，如果接收到的文本是HTML时执行
OnXML OnResponse 正确执行后，如果接收到的文本是XML时执行
OnScraped OnXML 回调后调用

代理
Colly通过其SetProxyFunc（）函数来切换代理。
任意自定义的函数可以通过SetProxyFunc()传参，只要这个函数签名为
	func(*http.Request) (*url.URL, error)

var proxies []*url.URL = []*url.URL{
    &url.URL{Host: "127.0.0.1:8080"},
    &url.URL{Host: "127.0.0.1:8081"},
}

func randomProxySwitcher(_ *http.Request) (*url.URL, error) {
    return proxies[random.Intn(len(proxies))], nil
}

// ...
c.SetProxyFunc(randomProxySwitcher)


Colly 有一个内置的代理切换器，可以在每个请求上轮流切换代理列表

Colly内置支持Google APP Engine。
如果你需要从标准的APP Engine中使用Colly，别忘了调用Collector.Appengine(*http.Request) 。
http://go-colly.org/docs/examples/scraper_server/

多收集器
http://go-colly.org/docs/examples/coursera_courses/

https://www.jianshu.com/p/fdf96da2d335
https://segmentfault.com/a/1190000019969473



https://www.xicidaili.com/wt/
https://www.kuaidaili.com/free/inha/
http://www.goubanjia.com/
http://go-colly.org/docs/introduction/start/


https://segmentfault.com/a/1190000020296085
https://segmentfault.com/a/1190000019969473
https://www.jianshu.com/p/23d4ecb8428f
https://www.codercto.com/a/32616.html
https://blog.csdn.net/sinat_36742186/article/details/85054317
*/
func main() {
	// 从Colly的repo中添加基本的日志调试器debug
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
	)
	// c.UserAgent = "xy"
	// c.AllowURLRevisit = true
	extensions.RandomUserAgent(c)

	if p, err := proxy.RoundRobinProxySwitcher(
		"socks5://127.0.0.1:1337",
		"socks5://127.0.0.1:1338",
		"http://127.0.0.1:8080",
	); err == nil {
		c.SetProxyFunc(p)
	}

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		fmt.Println("First column of a table row:", e.Text)
	})
	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	c.Visit("http://go-colly.org/")

}
