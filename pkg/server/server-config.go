package server

type ServiceConfig struct {
	Port         string   `json:"port" yaml:"port" toml:"port"`
	Address      string   `json:"address" yaml:"address" toml:"address"`
	ReadTimeout  int      `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout int      `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	Mode         string   `json:"mode" yaml:"mode" toml:"mode"`
	Middlewares  []string `json:"middlewares" yaml:"middlewares" toml:"middlewares"`
}
