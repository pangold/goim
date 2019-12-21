package interfaces

type Connection interface {
	GetToken() string
	Stop()
	Send([]byte)
	SetMessageHandler(*func([]byte, string) error)
}

type Pool interface {
	// for upstream invokers
	SetConnectedHandler(func(Connection) error)
	SetDisconnectedHandler(func(Connection) error)
	GetConnections() *map[string]Connection
	// for the same level invokers
	IsExist(Connection) bool
	Register(Connection)
	Unregister(Connection)
	NewConnection(Connection)
	LostConnection(Connection)
	HandleLoop()
}

type Server interface {
	GetPool() Pool
	Run()
}