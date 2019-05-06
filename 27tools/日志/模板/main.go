package main
import (
	"net/http"
	"html/template"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
)
func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	// router.Delims("{[{", "}]}")
	// 必须放在loadtemplate前面
	router.SetFuncMap(template.FuncMap{
        "formatAsDate": formatAsDate,
    })
	// router.LoadHTMLGlob("templates/**/*")
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")

	router.GET("index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
		})
	})
	router.GET("content-wrapper", func(c *gin.Context) {
		c.HTML(http.StatusOK, "content-wrapper.html", gin.H{
		})
	})
	router.GET("content-wrapper1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "content-wrapper.1.html", gin.H{
		})
	})
	router.GET("content-wrapper2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "content-wrapper.2.html", gin.H{
		})
	})
	router.GET("starter", func(c *gin.Context) {
		c.HTML(http.StatusOK, "starter.html", gin.H{
		})
	})

	router.Run(":8080")
}
func formatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d%02d/%02d", year, month, day)
}
// var tmp1 = "tmpl/tmpl.html"
// var tmp2 = "tmpl/tmpl2.html"
// func process(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("X-XSS-Protection","0")
// 	// t:=template.Must(template.ParseFiles("tmpl/tmpl.html"))
// 	t:=template.Must(template.ParseGlob("tmpl/*.html"))
// 	// t.Execute(w,"hello world")
// 	t.ExecuteTemplate(w,"layout","12333")
// }
// func main(){
// 	s:=http.Server{
// 		Addr:":8080",
// 	}
// 	http.HandleFunc("/process",process)
// 	s.ListenAndServe()
// }