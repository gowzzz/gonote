# 大文件读取和大请求写入

## 步骤
首先准备个服务



    // 把request的内容读取出来
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
    }
    // 把刚刚读出来的再写进去
    c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))


    func saveImag(sendUrl string, index int, imageChannel chan int)  {
	//创建文件
	path := "C:/img/" + strconv.Itoa(index) + ".jpg"
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	//获取http流
	resp, err :=http.Get(sendUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//分片逐步写入
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			break
		}
		f.Write(buf[:n])
	}
	imageChannel <- index
}
--------------------- 
作者：Hello_Ray 
来源：CSDN 
原文：https://blog.csdn.net/hello_ray/article/details/93332867 
版权声明：本文为博主原创文章，转载请附上博文链接！

package main

import (
    "log"
    "net/http"
    "os"
    "time"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
    video, err := os.Open("./test.mp4")
    if err != nil {
        log.Fatal(err)
    }
    defer video.Close()

    http.ServeContent(w, r, "test.mp4", time.Now(), video)
}

func main() {
    http.HandleFunc("/", ServeHTTP)
    http.ListenAndServe(":8080", nil)
}
https://blog.csdn.net/qq_30505673/article/details/90722014


func polling(ctx context.Context, incoming chan []byte) {
    http.HandleFunc("/polling", func(w http.ResponseWriter, r *http.Request) {
        select {
        case <- ctx.Done():
            fmt.Println("system quit")
        case b := <- incoming:
            w.Write(b)
        case <-time.After(5 * time.Second):
            w.Write("keepalive")
        case <- w.(http.CloseNotifier).CloseNotify():
            fmt.Println("connection closed")
        }
    })
}
--------------------- 
作者：win_lin 
来源：CSDN 
原文：https://blog.csdn.net/win_lin/article/details/78602130 
版权声明：本文为博主原创文章，转载请附上博文链接！