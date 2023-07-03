package server

import (
	"context"
	"os"
)

type Option func(o *options)

type options struct {
	id      string
	version string
	name    string
	sigs    []os.Signal
	ctx     context.Context

	handleSignals bool
	servers       []IAppServer
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Server with transport servers.
func Server(srv ...IAppServer) Option {
	return func(o *options) { o.servers = srv }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

func HandleSignals(handleSignals bool) Option {
	return func(o *options) { o.handleSignals = handleSignals }
}
