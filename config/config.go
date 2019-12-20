package config

type WsConfig struct {
	Protocol string
	Address  string
	CertFile string
	KeyFile  string
}

type TcpConfig struct {
	Addr     string
	Port     int
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
		Protocol: "ws",
		Ws: WsConfig{
			Protocol: "ws",
			Address:  "0.0.0.0:10000",
			CertFile: "",
			KeyFile: "",
		},
		Tcp: TcpConfig{
			Addr: "0.0.0.0",
			Port: 10001,
		},
	}
}