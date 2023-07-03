package config

import (
	"errors"

	"github.com/spf13/pflag"
)

type FlagSet = pflag.FlagSet

var ErrSkipProvider = errors.New("skip provider")

type (
	RegisterFlags func(flag *FlagSet)
	Unmarshaler   func(p []byte, v interface{}) error
	FlagParser    func() FlagParseResult
)

// Logger is the interface that wraps the basic Log methods for usage only in the config init stage
type Logger interface {
	// nolint: gofumpt
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

type options struct {
	usage            string
	shortDescription string
	flagParse        FlagParser
	logger           Logger
	registerFlags    RegisterFlags
	providers        []Provider // file, nacos, text
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithUsage(opt string) Option {
	return optionFunc(func(o *options) {
		o.usage = opt
	})
}

func WithShortDescription(opt string) Option {
	return optionFunc(func(o *options) {
		o.shortDescription = opt
	})
}

// WithRegisterFlags is used for add extra cli flags before parse args
func WithRegisterFlags(opt RegisterFlags) Option {
	return optionFunc(func(o *options) {
		o.registerFlags = opt
	})
}

func WithFlagParser(opt FlagParser) Option {
	return optionFunc(func(o *options) {
		o.flagParse = opt
	})
}

func WithProviders(opt ...Provider) Option {
	return optionFunc(func(o *options) {
		o.providers = opt
	})
}

// WithLogger the name ref to "Early KMS start"
func WithLogger(opt Logger) Option {
	return optionFunc(func(o *options) {
		o.logger = opt
	})
}
