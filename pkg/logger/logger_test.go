package logger

import "testing"

func show() {
	Info("info message")
	Infow("info message", "da", "123")
}

func TestLevel(t *testing.T) {
	show()
}

func TestColorLogger(t *testing.T) {
	config := NewDevelopmentConfig()
	config.Level = DebugLevel
	config.DisableStacktrace = true
	config.DisableCaller = true
	SetConfig(config)

	show()
}

func TestProdLogger(t *testing.T) {
	config := NewProductionConfig()
	config.Level = DebugLevel
	config.DisableStacktrace = true
	config.DisableCaller = true
	SetConfig(config)
	var l Level
	l.Set("debug")
	SetLevel(l)

	With("trace_id", "274ac2bbf9d5")
	With("span_id", "383d60f1")

	show()
}

func TestProdLoggerMap(t *testing.T) {
	config := NewProductionConfig(FieldPair{"service", "client_string"})
	config.Level = DebugLevel
	config.DisableStacktrace = true
	config.DisableCaller = true
	SetConfig(config)
	var l Level
	l.Set("debug")
	SetLevel(l)

	With("trace_id", "274ac2bbf9d5")
	With("span_id", "383d60f1")

	show()
}

func TestNewDefaultConfig(t *testing.T) {
	lc := NewDefaultConfig()
	lc.DisableStacktrace = true
	lc.EnableColor = true
	lc.Encoding = "console"
	lc.InitialFields = map[string]interface{}{
		"service": "client_string",
	}

	SetConfig(lc)

	show()
}
