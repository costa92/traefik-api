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
