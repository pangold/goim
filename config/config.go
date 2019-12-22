package config

type WsConfig struct {
	Protocol  string
	Address   string
	CertFile  string
	KeyFile   string
}

type TcpConfig struct {
	Address   string
}

type Config struct {
	Protocols []string    // support multi protocols
	Ws        WsConfig    // config of websocket protocol
	Wss       WsConfig    // https
	Tcp       TcpConfig   // config of tcp protocol
}

var (
	Conf Config
)

func init() {
	Conf = Config{
		Protocols: []string{"tcp", "ws"},
		Ws: WsConfig{
			Protocol: "ws",
			Address:  "0.0.0.0:10000",
		},
		Wss: WsConfig{
			Protocol: "wss",
			Address:  "0.0.0.0:10001",
			CertFile: "",
			KeyFile:  "",
		},
		Tcp: TcpConfig{
			Address:  "127.0.0.1:10002",
		},
	}
}