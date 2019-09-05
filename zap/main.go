package main

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "fmt"
    "time"
)

func main() {
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
        EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
    }

    // 设置日志级别
    atom := zap.NewAtomicLevelAt(zap.DebugLevel)

    config := zap.Config{
        Level:            atom,                                                // 日志级别
        Development:      true,                                                // 开发模式，堆栈跟踪
        Encoding:         "json",                                              // 输出格式 console 或 json
        EncoderConfig:    encoderConfig,                                       // 编码器配置
        InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
        OutputPaths:      []string{"stdout", "./logs/spikeProxy.log"},         // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
        ErrorOutputPaths: []string{"stderr"},
    }

    // 构建日志
    logger, err := config.Build()
    if err != nil {
        panic(fmt.Sprintf("log 初始化失败: %v", err))
    }
    logger.Info("log 初始化成功")

    logger.Info("无法获取网址",
        zap.String("url", "http://www.baidu.com"),
        zap.Int("attempt", 3),
        zap.Duration("backoff", time.Second),
    )
}
// #输出结果
// #{"level":"info","time":"2019-01-02T15:38:33.778+0800","caller":"/Users/lcl/go/src/spikeProxy/main.go:54","msg":"log 初始化成功","serviceName":"spikeProxy"}
// #{"level":"info","time":"2019-01-02T15:38:33.778+0800","caller":"/Users/lcl/go/src/spikeProxy/main.go:56","msg":"无法获取网址","serviceName":"spikeProxy","url":"http://www.baidu.com","attempt":3,"backoff":1}