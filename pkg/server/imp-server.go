package server

import "context"

type IAppServer interface {
	Start(context.Context) error
	Stop(context.Context) error
}
