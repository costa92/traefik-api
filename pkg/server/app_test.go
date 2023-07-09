package server

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app := New()
	err := app.Run()
	if err != nil {
		t.Fatal(err)
	}
}
