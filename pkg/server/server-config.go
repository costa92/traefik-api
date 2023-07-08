package server

type ServiceConfig struct {
	Name         string   `json:"name" yaml:"name" toml:"name"`
	Port         string   `json:"port" yaml:"port" toml:"port" json:"port,omitempty"`
	Address      string   `json:"address" yaml:"address" toml:"address" json:"address,omitempty"`
	ReadTimeout  int      `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" json:"readTimeout,omitempty"`
	WriteTimeout int      `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" json:"writeTimeout,omitempty"`
	Mode         string   `json:"mode" yaml:"mode" toml:"mode" json:"mode,omitempty"`
	Middlewares  []string `json:"middlewares" yaml:"middlewares" toml:"middlewares" json:"middlewares,omitempty"`
}
