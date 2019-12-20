package config

type WsConfig struct {
	Protocol string
	Address  string
	CertFile string
	KeyFile  string
}

type TcpConfig struct {
	Address  string
}

type Config struct {
	Protocol string
	Ws       WsConfig
	Tcp      TcpConfig
}

var (
	Conf Config
)

func init() {
	Conf = Config{
		Protocol: "tcp",
		Ws: WsConfig{
			Protocol: "ws",
			Address:  "0.0.0.0:10000",
			CertFile: "",
			KeyFile:  "",
		},
		Tcp: TcpConfig{
			Address:  "127.0.0.1:10001",
		},
	}
}