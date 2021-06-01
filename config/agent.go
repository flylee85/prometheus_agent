package config

type AgentConfig struct {
	UUID         string
	Addr         string
	ServerConfig `mapstructure:"server"`
	LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Addr  string
	Token string
}

type LogConfig struct {
	Filename   string
	Maxsize    int
	Maxbackups int
	Compress   bool
}
