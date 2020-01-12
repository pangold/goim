package config

type WsConfig struct {
	Protocol  string       `yaml:"protocol"`
	Address   string       `yaml:"address"`
	CertFile  string       `yaml:"cert"`
	KeyFile   string       `yaml:"key"`
}

type HostConfig struct {
	Address   string       `yaml:"address"`
}

type Token struct {
	SecretKey string       `yaml:"secret-key"`
}

// For frontend server/client
type FrontConfig struct {
	Protocols []string     `yaml:"protocols"`
	Ws        WsConfig     `yaml:"ws"`
	Wss       WsConfig     `yaml:"wss"`
	Tcp       HostConfig   `yaml:"tcp"`
}

// For backend server
type BackConfig struct {
	Http      HostConfig   `yaml:"http"`
	Grpc      HostConfig   `yaml:"grpc"`
}

type Config struct {
	Front     FrontConfig  `yaml:"front"`
	Back      BackConfig   `yaml:"back"`
	Token     Token        `yaml:"token"`
}