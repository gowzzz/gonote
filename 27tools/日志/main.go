package main

import (
  "github.com/sirupsen/logrus"
  "os"
  "runtime"
	//   "fmt"
	//   "strings"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"time"
)

func main() {
// logger是一种相对高级的用法, 对于一个大型项目, 往往需要一个全局的logrus实例，即logger对象来记录项目所有的日志。
	var log = logrus.New()
	log.SetReportCaller(true)
	 // 设置日志级别为xx以及以上
	 log.SetLevel(logrus.InfoLevel)
	log.AddHook(&DefaultFieldHook{})
	// 设置日志格式为json格式
	// log.SetFormatter(&logrus.JSONFormatter{
	// 	// PrettyPrint: true,//格式化json
	// 	TimestampFormat: "2006-01-02 15:04:05",//时间格式化
	// })
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:true,
		EnvironmentOverrideColors:true,
		// FullTimestamp:true,
		TimestampFormat: "2006-01-02 15:04:05",//时间格式化
		// DisableLevelTruncation:true,
	})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)
	// // 缩短文件名
	// log.Formatter = &logrus.TextFormatter{
	// 	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	// 		repopath := fmt.Sprintf("%s/src/github.com/bob", os.Getenv("GOPATH"))
	// 		filename := strings.Replace(f.File, repopath, "", -1)
	// 		return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
	// 	},
	// }


	logName:=`G:\gopath\src\github.com\gonote\27tools\日志\aaa`
	writer, err := rotatelogs.New(
        logName+".%Y%m%d%H%M%S",
        // WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
        // rotatelogs.WithLinkName(logName),

        // WithRotationTime 设置日志分割的时间，这里设置为一小时分割一次
        rotatelogs.WithRotationTime(time.Minute*1),//文件旋转之间的间隔。默认情况下，日志每86400秒/一天旋转一次。注意:记住要利用时间。持续时间值。
        // WithMaxAge和WithRotationCount二者只能设置一个，
        // WithMaxAge设置文件清理前的最长保存时间，
        // WithRotationCount设置文件清理前最多保存的个数。 默认情况下，此选项是禁用的。
		// rotatelogs.WithMaxAge(time.Second*30),//默认每7天清除下日志文件
		rotatelogs.WithMaxAge(-1), //需要手动禁用禁用  默认情况下不清除日志，
        // rotatelogs.WithRotationCount(2),//清除除最新2个文件之外的日志，默认禁用
    )
    if err != nil {
        log.Errorf("config local file system for logger error: %v", err)
    }

    lfsHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer,
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
    }, &logrus.TextFormatter{DisableColors: true})

	log.AddHook(lfsHook)



	// 初始化一些公共参数
	loginit:=log.WithFields(logrus.Fields{
		"animal": "walrus",
	})

	for i:=0;i<10000;i++{
		go func(log *logrus.Entry){
			for {
				time.Sleep(time.Second*1)
				log.Info("A walrus appears")
			}
		}(loginit)
	
	}

	select{}
}
type DefaultFieldHook struct {
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	funcName, file, line, ok := runtime.Caller(0)
	if ok{
		entry.Data["funcName"] = runtime.FuncForPC(funcName).Name()
		entry.Data["file"] = file
		entry.Data["line"] = line
	}

		entry.Data["appName"] = "appName"
   
    return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
    return logrus.AllLevels
}


