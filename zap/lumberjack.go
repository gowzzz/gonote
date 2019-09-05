package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "time"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
)

func main() {
    hook := lumberjack.Logger{
        Filename:   "./logs/spikeProxy1.log", // 日志文件路径
        MaxSize:    128,                      // 每个日志文件保存的最大尺寸 单位：M
        MaxBackups: 30,                       // 日志文件最多保存多少个备份
        MaxAge:     7,                        // 文件最多保存多少天
        Compress:   true,                     // 是否压缩
    }

    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "linenum",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
        EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
        EncodeDuration: zapcore.SecondsDurationEncoder, //
        EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
        EncodeName:     zapcore.FullNameEncoder,
    }

    // 设置日志级别
    atomicLevel := zap.NewAtomicLevel()
    atomicLevel.SetLevel(zap.ErrorLevel)
	/*
	
	DebugLevel = zapcore.DebugLevel
    InfoLevel = zapcore.InfoLevel
    WarnLevel = zapcore.WarnLevel
    ErrorLevel = zapcore.ErrorLevel
    DPanicLevel = zapcore.DPanicLevel
    PanicLevel = zapcore.PanicLevel
    FatalLevel = zapcore.FatalLevel
	*/ 
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
        zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
        atomicLevel,                                                                     // 日志级别
    )

    // 开启开发模式，堆栈跟踪
    caller := zap.AddCaller()
    // 开启文件及行号
    development := zap.Development()
    // 设置初始化字段
    filed := zap.Fields(zap.String("serviceName", "serviceName"))
    // 构造日志
    logger := zap.New(core, caller, development, filed)

    logger.Debug("log DebugDebugDebug")
    logger.Info("log InfoInfoInfoInfo")
    logger.Warn("log WarnWarnWarnWarn")
    logger.Error("log ErrorErrorErrorError")
    // logger.DPanic("log DPanicDPanicDPanicDPanic")//会退出
    // logger.Panic("log PanicPanicPanicPanic")/会退出
    // logger.Fatal("log FatalFatalFatalFatalFatal")//会退出
    logger.Info("Info InfoInfo",
        zap.String("url", "http://www.baidu.com"),
        zap.Int("attempt", 3),
        zap.Duration("backoff", time.Second))
}
// # 控制台输出结果
// # {"level":"info","time":"2019-01-02T16:14:43.608+0800","linenum":"/Users/lcl/go/src/spikeProxy/main.go:56","msg":"log 初始化成功","serviceName":"serviceName"}
// # {"level":"info","time":"2019-01-02T16:14:43.608+0800","linenum":"/Users/lcl/go/src/spikeProxy/main.go:57","msg":"无法获取网址","serviceName":"serviceName","url":"http://www.baidu.com","attempt":3,"backoff":1}
// # 文件输出结果
// # {"level":"info","time":"2019-01-02T16:14:43.608+0800","linenum":"/Users/lcl/go/src/spikeProxy/main.go:56","msg":"log 初始化成功","serviceName":"serviceName"}
// # {"level":"info","time":"2019-01-02T16:14:43.608+0800","linenum":"/Users/lcl/go/src/spikeProxy/main.go:57","msg":"无法获取网址","serviceName":"serviceName","url":"http://www.baidu.com","attempt":3,"backoff":1}