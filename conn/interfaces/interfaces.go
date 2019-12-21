package interfaces

// Conn is a generic connection IM connection
type Conn interface {
	// Get token of the this connection
	GetToken() string
	// Close connection immediately and directly
	Close()
	// Stop sending and receiving loop.
	// Close connection after loops exit completely.
	Stop()
	// Send message to this connection asynchronisely.
	Send([]byte)
	// Received message handler for expansibilities.
	// How do you want to deal with the received data?
	// Parsing it? Wanted to ignore it?
	// Transferring to others?
	// Do whatever you wanted to in passing handler.
	SetMessageHandler(*func([]byte, string) error)
}

// Pool is a generic connection pool that manages generic Conn,
//
// New detected and undetected connections will be callback
// if specifics handlers to handle some special extra works.
// For example, for session management or synchronizing session.
type Pool interface {
	// For downstream invokers.

	// New connection callback that tells invoker you have a visitor.
	// Mark it down in session tables?
	// Pass it to DB/Cache Server?
	// Or just ignore it if your don't want it.
	SetConnectedHandler(func(Conn) error)
	// New disconnection callback that tells invoker that a visitor left.
	// Erase it from session tables?
	// Pass it to DB/Cache Server to erase it?
	// Or just ignore it if your don't want it.
	SetDisconnectedHandler(func(Conn) error)
	// Get all connections that is currently connected
	// and be managed in connection pool.
	GetConnections() *map[string]Conn

	// For the same level invokers: TCP and Websocket

	// Check if a connection is still connected.
	IsExist(Conn) bool
	// A new connection register.
	// To manage concurrency control.
	Register(Conn)
	// A new disconnection.
	// To manage concurrency control.
	Unregister(Conn)
	// Manage a loop dealing with new connections and disconnections.
	HandleLoop()
}

// A server that takes control of all components(Pool, Connections)
// to provide services.
type Server interface {
	// Manage a loop handling with all behaviors of connections.
	Run()
	// Get connection pool.
	// Basically for communication with connections.
	GetPool() Pool
}