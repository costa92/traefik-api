package server

import (
	"context"
	"net/http"
)

type IAppServer interface {
	PreRun(router http.Handler)
	Start(context.Context) error
	Stop(context.Context) error
}
