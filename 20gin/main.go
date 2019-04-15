package main

import (
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

/*
GET("/user/:name/*action", func  =>  c.Param("name")   c.Param("action")
/welcome?firstname=Jane&lastname=Doe => c.DefaultQuery("firstname", "Guest")  c.Query("lastname")
Multipart/Urlencoded Form  =>c.PostForm("message")   c.DefaultPostForm("nick", "anonymous")
post  names[first]=thinkerou&names[second]=tianou  对象？？  => c.PostFormMap("names")
?ids[a]=1234&ids[b]=hello HTTP/1.1  => c.QueryMap("ids")
组：v1 := router.Group("/v1")
{

}

gin.Default() =>没有中间件的空白Gin
*/
type Login struct {
	User     string    `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string    `form:"password" json:"password" xml:"password" binding:"required"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
	// ID       string    `uri:"id" binding:"required,uuid"`
}

/*
ShouldBindJSON: json(raw)
ShouldBind: json xml form-data param x-www-form-irlencoded 当param和其他同时存在时以其他的为最终数据
ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
ShouldBindWith(&json, binding.Query)  =>  xx?a=1&b=2
c.BindQuery function that only binds the query params and not the post data.
time_format:"2006-01-02" time_utc:"1"

uri：ID string `uri:"id" binding:"required,uuid"`   c.ShouldBindUri(&person)
*/

// Booking contains binded and validated data.  ltfield 小于  gtfield
type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,ltfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}

func main() {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	// r := gin.New()
	// Recovery返回一个中间件，该中间件从任何恐慌中恢复，如果有500，则写入500。
	v1 := r.Group("/v1")
	v1.Use(gin.Logger())
	v1.Use(Logger())
	{
		v1.POST("/ping/:id", func(c *gin.Context) {
			var json Login
			if err := c.ShouldBind(&json); err != nil {
				c.SecureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else {
				c.SecureJSON(http.StatusOK, json)
			}
		})
		v1.POST("/book", func(c *gin.Context) {
			var b Booking
			if err := c.ShouldBind(&b); err != nil {
				c.SecureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			} else {
				c.SecureJSON(http.StatusOK, b)
			}
		})
	}

	http.ListenAndServe(":8888", r)
	r.Run(":8888") // listen and serve on 0.0.0.0:8080
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

// 一般为组添加中间件

/*
gin专门的middleware中间件仓库：contrib  用来存储大家贡献的中间件
适合在中间件中做的事情
1.对response body进行压缩
2. 设置特殊路由
3. 打印request处理日志
4.挂在proff需要的路由
5.将http request中的remoteaddr改成realip
6.为本次请求产生单独的requestid
7.设置context.timeout的超时时间
8.存token，对接口限流

验证一般和结构体绑定

*/

/*
模型绑定和验证
We currently support binding of JSON, XML, YAML and standard form values (foo=bar&boo=baz).
gin use this：https://github.com/go-playground/validator

gin提供两套绑定方法：
Must bind：这个如果出错会直接返回400，你不能再次修改错误码。如果你想这么做用下面的should bind
	Bind, BindJSON, BindXML, BindQuery, BindYAML
Should bind：如果有问题，返回错误，开发人员有责任适当地处理请求和错误
	ShouldBind, ShouldBindJSON, ShouldBindXML, ShouldBindQuery, ShouldBindYAML
还可以使用binding来限制非空，会直接报错

还可以注册自定义验证器。在binding后面添加自定义的方法名
ShouldBindQuery函数只绑定查询参数，而不绑定post数据
c.ShouldBind(&person)=》Bind Query String or Post Data
*/
