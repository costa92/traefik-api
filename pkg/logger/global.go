package logger

var (
	defaultConfig = NewDefaultConfig()
	l             = newLogger(defaultConfig)
)

func SetConfig(config *Config) {
	l = newLogger(config)
}

func GetConfig() *Config {
	return l.config
}

func SetLevel(level Level) {
	l.SetLevel(level)
}

func GetLevel() Level {
	return l.Level()
}

func SetOutputPaths(outputPaths []string) {
	l.config.OutputPaths = outputPaths
	l = newLogger(l.config)
}

func NewLogger(config *Config) Logger {
	l := newLogger(config)
	return l
}

func GetLogger() Logger {
	return l
}
