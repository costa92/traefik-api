package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/costa92/errors"
	"github.com/segmentio/ksuid"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"

	"treafik-api/pkg/logger"
)

// AppInfo is application context value.
type AppInfo interface {
	ID() string
	Name() string
	Version() string
}

type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

func New(opts ...Option) *App {
	options := options{
		ctx:           context.Background(),
		sigs:          []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		handleSignals: true,
	}
	options.id = ksuid.New().String()
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	// nolint: forbidigo
	if _, err := maxprocs.Set(maxprocs.Logger(logger.Infof)); err != nil {
		return err
	}

	ctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()
	c := make(chan os.Signal, 1)

	// warning: you need manually call App.Stop() to stop the application if you set handleSignals to false
	if a.opts.handleSignals {
		signal.Notify(c, a.opts.sigs...)
	}

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return a.Stop()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
