package gin_server

import (
	"time"

	"treafik-api/pkg/server"
)

var _ server.IAppServer = (*AppServer)(nil)

type ServerOption func(appServer *AppServer)

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *AppServer) {
		s.timeout = timeout
	}
}

func WithMiddleware(middlewares []string) ServerOption {
	return func(s *AppServer) {
		s.middleware = middlewares
	}
}
