// kfhjashvuvxrcabj
// package main

// import (
// 	// "gopkg.in/gomail.v2"
// )
package main

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

func main() {
	fromUser := "golang<371600645@qq.com>"
	toUser := "810169879@qq.com"
	subject := "hello,world"
	// NewEmail返回一个email结构体的指针
	e := email.NewEmail()
	// 发件人
	e.From = fromUser
	// 收件人(可以有多个)
	e.To = []string{toUser}
	// 邮件主题
	e.Subject = subject
	// // 解析html模板
	// t, err := template.ParseFiles("email-template.html")
	// if err != nil {
	// 	return err
	// }
	// // Buffer是一个实现了读写方法的可变大小的字节缓冲
	// body := new(bytes.Buffer)
	// // Execute方法将解析好的模板应用到匿名结构体上，并将输出写入body中
	// t.Execute(body, struct {
	// 	FromUserName string
	// 	ToUserName   string
	// 	TimeDate     string
	// 	Message      string
	// }{
	// 	FromUserName: "go语言",
	// 	ToUserName:   "Sixah",
	// 	TimeDate:     time.Now().Format("2006/01/02"),
	// 	Message:      "golang是世界上最好的语言！",
	// })
	// // html形式的消息
	// e.HTML = body.Bytes()
	// 从缓冲中将内容作为附件到邮件中
	// e.Attach(body, "email-template.html", "text/html")
	// e.AttachFile("/home/shuai/go/src/email/main.go")
	// 发送邮件(如果使用QQ邮箱发送邮件的话，passwd不是邮箱密码而是授权码)
	e.Send("smtp.qq.com:587", smtp.PlainAuth("用户名", "371600645@qq.com", "kfhjashvuvxrcabj", "smtp.qq.com"))
}

// func main() {
// 	m := gomail.NewMessage()

// 	m.SetHeader("From", "371600645@qq.com")
// 	m.SetHeader("To", "810169879@qq.com", "13521146683@163.com")
// 	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
// 	m.SetHeader("Subject", "Hello!")
// 	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
// 	//m.Attach("/home/Alex/lolcat.jpg")
// 	m.Attach()
// 	m.AddAlternative()
// 	m.AddAlternativeWriter()
// 	d := gomail.NewDialer("smtp.qq.com", 587, "371600645@qq.com", "kfhjashvuvxrcabj")

// 	// Send the email to Bob, Cora and Dan.
// 	if err := d.DialAndSend(m); err != nil {
// 		panic(err)
// 	}
// }
