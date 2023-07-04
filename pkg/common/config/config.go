package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"honnef.co/go/tools/version"
)

var ErrEmptyConfig = errors.New("get empty config")

type FlagParseResult interface {
	ConfigFile() string
	DumpConfig() bool
	ShowHelp() bool
	ShowVersion() bool
	Usage() func()
}

const (
	EnvNacosHost      = "NACOS_HOST"
	EnvNacosNamespace = "NACOS_NAMESPACE"
	EnvNacosPort      = "NACOS_PORT"
	EnvNacosGroup     = "NACOS_GROUP"
	EnvNacosDataID    = "NACOS_DATAID"
	EnvNacosLogLevel  = "NACOS_LOG_LEVEL"

	FlagConfigFile = "config"
	FlagDumpConfig = "dump"
)

type ConfigLoader struct {
	options    options
	configFile string
}

func New(opts ...Option) *ConfigLoader {
	options := options{
		usage:            fmt.Sprintf("Usage: %s [Options]", "v1"),
		shortDescription: fmt.Sprintf("%s %s", "ServiceName", "info"),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	return &ConfigLoader{
		options: options,
	}
}

func (cl *ConfigLoader) Load(cfg interface{}) error {
	mtype := reflect.TypeOf(cfg)
	if mtype.Kind() != reflect.Ptr {
		return errors.New("only a pointer to struct or map can be unmarshalled from config content")
	}
	var flagResult FlagParseResult
	if cl.options.flagParse != nil {
		flagResult = cl.options.flagParse()
	} else {
		flagResult = cl.defaultFlagParser()
	}
	cl.configFile = flagResult.ConfigFile()

	if flagResult.ShowHelp() {
		flagResult.Usage()()
		os.Exit(0)
	}

	if flagResult.ShowVersion() {
		fmt.Println("v1")
		os.Exit(0)
	}
	// 日志
	if cl.options.logger == nil {
		zapCfg := zap.NewProductionConfig()
		zapCfg.InitialFields = map[string]interface{}{
			"service": "ServiceName",
			"version": version.Version,
		}
		zapCfg.OutputPaths = []string{"stderr"}

		zapLogger, err := zapCfg.Build()
		if err != nil {
			panic(err)
		}
		cl.options.logger = zapLogger.Sugar()
	}

	err := cl.getConfigViaProviders(cfg)
	if err != nil {
		return err
	}
	return nil
}

type defaultFlagResult struct {
	configFile  string
	dumpConfig  bool
	showHelp    bool
	showVersion bool
	usage       func()
}

func (f *defaultFlagResult) ConfigFile() string {
	return f.configFile
}

func (f *defaultFlagResult) DumpConfig() bool {
	return f.dumpConfig
}

func (f *defaultFlagResult) ShowHelp() bool {
	return f.showHelp
}

func (f *defaultFlagResult) ShowVersion() bool {
	return f.showVersion
}

func (f *defaultFlagResult) Usage() func() {
	return f.usage
}

func (cl *ConfigLoader) defaultFlagParser() FlagParseResult {
	var configFile string
	var dumpConfig bool
	var showHelp, showVersion bool
	commandLine := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	// use standalone instead of shared default pflag.CommandLine avoid "pflag redefined: config" error when unit tests
	commandLine.Usage = func() {
		fmt.Fprint(os.Stderr, cl.options.usage, "\n\n")
		fmt.Fprint(os.Stderr, cl.options.shortDescription, "\n\n")
		fmt.Fprintln(os.Stderr, commandLine.FlagUsages())
	}
	commandLine.SortFlags = false

	commandLine.StringVarP(&configFile, FlagConfigFile, "c", "", "config file path")
	commandLine.BoolVar(&dumpConfig, FlagDumpConfig, false, "dump config to toml")
	commandLine.BoolVarP(&showVersion, "version", "v", false, "display the current version of this CLI")
	commandLine.BoolVarP(&showHelp, "help", "h", false, "show help")

	if cl.options.registerFlags != nil {
		cl.options.registerFlags(commandLine)
	}
	commandLine.Parse(os.Args[1:])
	return &defaultFlagResult{
		configFile,
		dumpConfig,
		showHelp,
		showVersion,
		commandLine.Usage,
	}
}

func (cl *ConfigLoader) getConfigViaProviders(cfg interface{}) error {
	var err error
	helper := &providerHelper{
		configFile: cl.configFile,
		log:        cl.options.logger,
	}

	for _, provider := range cl.options.providers {
		_, err = provider.Config(helper, cfg)
		if err == nil {
			break
		}
		if errors.Is(err, ErrSkipProvider) {
			cl.options.logger.Infow("config provider skipped", "provider", provider.Name(), "reason", err)
		} else {
			cl.options.logger.Infow("try get config via provider failed", "provider", provider.Name(), "err", err)
		}
	}
	if len(cl.options.providers) == 0 {
		err = errors.New("error no config provider usable")
	}
	return err
}
