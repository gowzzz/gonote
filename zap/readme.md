# zap.NewDevelopment() 简单示例，格式化输出：
2018-01-18T15:40:05.991+0800    INFO    tool/zaplog.go:83       to sugar failed to fetch URLurlhttp://example.comattempt3backoff1s

# zap.NewProduction() json序列化输出
{"level":"info","ts":1516261205.991458,"caller":"tool/zaplog.go:109","msg":"to sugar failed to fetch URLurlhttp://example.comattempt3backoff1s"}

# 动态改变日志的打印级别
core := zapcore.NewCore(
    zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), //开发者Encoder，包含函数调用信息
    // zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
    w,
    // zap.InfoLevel, //info,error
    atom, //debug,info,warn,error
)
logger := zap.New(core)

---------------
atom := zap.NewAtomicLevel()
atom.SetLevel(zap.DebugLevel)

{"L":"INFO","T":"2018-01-18T16:18:05.324+0800","M":"to desugar failed to fetch URL","url":"http://example.com","attempt":3,"backoff":"1s"}



