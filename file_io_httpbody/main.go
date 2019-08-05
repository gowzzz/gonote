package main 

import(
	"bytes"
	"net/http"
	"os"
	"log"
	"io"
	"fmt"
	"net/url"

)
func test1(w http.ResponseWriter, r *http.Request){
	var buf bytes.Buffer
	buf.WriteString("1111111")
	w.Write(buf.Bytes())
}
func test2(w http.ResponseWriter, r *http.Request){
	v := url.Values{}
	v.Add("filename", "2.mkv")
	w.Header().Set("Content-Disposition", "attachment;"+v.Encode())
	video, err := os.Open("./1.mkv")
    if err != nil {
        log.Fatal(err)
    }
	defer video.Close()
	num:=0
	buf := make([]byte, 1)//1M
	for {
		select {
		case <- w.(http.CloseNotifier).CloseNotify():
			fmt.Println("connection closed")
			break
		default:
			num++
			n, err := video.Read(buf)
			fmt.Println("n:",n)
			fmt.Println("num:",num)
			w.Write(buf[:n])
			if err == io.EOF {
				fmt.Println("err:",err)
				break
			}
			if err != nil {
				fmt.Println("err2:",err)
				break
			}
		}
		
	}
}
func main(){
	server:=http.Server{
		Addr:":9999",
	}
	http.HandleFunc("/test1",test1)
	http.HandleFunc("/test2",test2)
	if err:=server.ListenAndServe();err!=nil{
		panic(err)
	}
}

func test(){
	video, err := os.Open("./1.mkv")
    if err != nil {
        log.Fatal(err)
    }
	defer video.Close()

	f, err := os.Create("./2.mkv")
    if err != nil {
        log.Fatal(err)
    }
	defer f.Close()
	num:=0
	buf := make([]byte, 1024*1024)//1M
	for {
		num++
		n, err := video.Read(buf)
		fmt.Println("n:",n)
		fmt.Println("num:",num)
		f.Write(buf[:n])
		if err == io.EOF {
			fmt.Println("err:",err)
			break
		}
		if err != nil {
			fmt.Println("err2:",err)
			break
		}
	}
}