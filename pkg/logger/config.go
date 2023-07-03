package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FieldPair []string

// zap v1.24.0
type Config struct {
	// 动态改变日志级别，在运行时你可以安全改变日志级别
	Level Level `json:"level" yaml:"level"`
	// 将日志记录器设置为开发模式，在 WarnLevel 及以上级别日志会包含堆栈跟踪信息
	Development bool `json:"development" yaml:"development"`
	// 在日志中停止调用函数所在文件名、行数
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// 完全禁止自动堆栈跟踪。默认情况下，在 development 中，warnlevel及以上日志级别会自动捕获堆栈跟踪信息
	// 在 production 中，ErrorLevel 及以上也会自动捕获堆栈信息
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// 设置日志编码。可以设置为 console 和 json。也可以通过 RegisterEncoder 设置第三方编码格式
	Encoding string `json:"encoding" yaml:"encoding"`
	// 为encoder编码器设置选项。详细设置信息在 zapcore.zapcore.EncoderConfig
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	// 日志输出地址可以是一个 URLs 地址或文件路径，可以设置多个
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	// 错误日志输出地址。默认输出标准错误信息
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	// 可以添加自定义的字段信息到 root logger 中。也就是每条日志都会携带这些字段信息，公共字段
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`

	EnableColor bool
	ShortTime   bool

	CallerSkip int
	zapConfig  *zap.Config
}

type ConfigInterface interface {
	GetLevel() string
	GetOutput() string
	GetEncoding() string
	GetDisableStacktrace() bool
	GetInitialFields() map[string]interface{}
}

func (c *Config) GetLevel() string {
	return c.Level.String()
}

func (c *Config) GetOutput() string {
	if len(c.OutputPaths) > 0 {
		return c.OutputPaths[0]
	}
	return "stderr"
}

func (c *Config) GetEncoding() string {
	return c.Encoding
}

func (c *Config) GetDisableStacktrace() bool {
	return c.DisableStacktrace
}

func NewProductionConfig(fields ...FieldPair) *Config {
	return &Config{
		Level:             InfoLevel,
		Development:       false,
		Encoding:          "json",
		OutputPaths:       []string{"stderr"},
		CallerSkip:        2,
		DisableStacktrace: false,
		InitialFields:     genInitialFields(fields),
	}
}

func NewDevelopmentConfig(fields ...FieldPair) *Config {
	return &Config{
		Level:             DebugLevel,
		Development:       true,
		ShortTime:         true,
		EnableColor:       true,
		Encoding:          "console",
		OutputPaths:       []string{"stderr"},
		CallerSkip:        2,
		DisableStacktrace: true,
		InitialFields:     genInitialFields(fields),
	}
}

func genInitialFields(args []FieldPair) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, pair := range args {
		fields[pair[0]] = pair[1]
	}
	return fields
}

func NewDefaultConfig(fields ...FieldPair) *Config {
	return NewProductionConfig(fields...)
}

func NewConfigFromInterface(config ConfigInterface) *Config {
	dft := NewDefaultConfig()

	if config.GetLevel() != "" {
		if lvl, err := ParseLevel(config.GetLevel()); err == nil {
			dft.Level = lvl
		}
	}

	if config.GetOutput() != "" {
		dft.OutputPaths = []string{config.GetOutput()}
	}

	if config.GetEncoding() != "" {
		dft.Encoding = config.GetEncoding()
	}
	dft.DisableStacktrace = config.GetDisableStacktrace()
	dft.InitialFields = config.GetInitialFields()
	return dft
}

func (c *Config) buildZapConfig() {
	encoderConfig := c.newCustomEncoderConfig()

	zapConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(c.Level)),
		Development:       c.Development,
		DisableCaller:     c.DisableCaller,
		DisableStacktrace: c.DisableStacktrace,
		Sampling:          &zap.SamplingConfig{Initial: 100, Thereafter: 100},
		Encoding:          c.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       c.OutputPaths,
		InitialFields:     c.InitialFields,
	}
	c.zapConfig = zapConfig
}

func (c *Config) newCustomEncoderConfig() zapcore.EncoderConfig {
	encoderLevel := zapcore.LowercaseLevelEncoder
	if c.EnableColor {
		encoderLevel = zapcore.LowercaseColorLevelEncoder
	}
	encodeTime := zapcore.ISO8601TimeEncoder
	if c.ShortTime {
		encodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			type appendTimeEncoder interface {
				AppendTimeLayout(time.Time, string)
			}
			layout := "2006-01-02 15:04:05"
			if enc, ok := enc.(appendTimeEncoder); ok {
				enc.AppendTimeLayout(t, layout)
				return
			}
			enc.AppendString(t.Format(layout))
		}
	}

	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encoderLevel,
		EncodeTime:     encodeTime,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (c *Config) clone() *Config {
	cloned := *c
	cloned.OutputPaths = make([]string, len(c.OutputPaths))
	copy(cloned.OutputPaths, c.OutputPaths)
	cloned.InitialFields = make(map[string]interface{})
	for k, v := range c.InitialFields {
		cloned.InitialFields[k] = v
	}
	if cloned.zapConfig != nil {
		cloned.buildZapConfig()
	}
	return &cloned
}
