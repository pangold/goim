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

type HttpConfig struct {
	Address   string
}

type GrpcConfig struct {
	Address   string
}

type Token struct {
	SecretKey string
}

type Config struct {
	Protocols []string    // support multi protocols
	// For frontend server/client
	Ws        WsConfig    // config of websocket protocol
	Wss       WsConfig    // https
	Tcp       TcpConfig   // config of tcp protocol
	// For backend server
	Http      HttpConfig
	Grpc      GrpcConfig
	Token     Token
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
			Address:  "0.0.0.0:10002",
		},
		Http: HttpConfig{
			Address:  "0.0.0.0:10003",
		},
		Grpc: GrpcConfig{
			Address:  "0.0.0.0:10004",
		},
		Token: Token{
			SecretKey:"my-secret-key",
		},
	}
}