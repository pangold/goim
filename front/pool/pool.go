package pool

import (
	"gitlab.com/pangold/goim/front/interfaces"
	"log"
)

type Pool struct {
	connections          map[string]interfaces.Conn
	connectedHandler    *func(interfaces.Conn) error
	disconnectedHandler *func(interfaces.Conn) error
	register             chan interfaces.Conn
	unregister           chan interfaces.Conn
}

func NewPool() interfaces.Pool {
	return &Pool{
		connections:         make(map[string]interfaces.Conn),
		connectedHandler:    nil,
		disconnectedHandler: nil,
		register:            make(chan interfaces.Conn),
		unregister:          make(chan interfaces.Conn),
	}
}

// for upstream invokers
func (p *Pool) SetConnectedHandler(fn func(interfaces.Conn) error) {
	p.connectedHandler = &fn
}

func (p *Pool) SetDisconnectedHandler(fn func(interfaces.Conn) error) {
	p.disconnectedHandler = &fn
}

func (p *Pool) GetConnections() *map[string]interfaces.Conn {
	return &p.connections
}

// for the same level invokers
func (p *Pool) IsExist(conn interfaces.Conn) bool {
	if _, ok := p.connections[conn.GetToken()]; ok {
		return true
	}
	return false
}

func (p *Pool) Register(conn interfaces.Conn) {
	p.register <- conn
}

func (p *Pool) Unregister(conn interfaces.Conn) {
	p.unregister <- conn
}

func (p *Pool) NewConnection(conn interfaces.Conn) {
	if p.connectedHandler != nil {
		// check if it is valid connection(actually, check for token)
		if err := (*p.connectedHandler)(conn); err != nil {
			log.Println("failed to verify, ", err.Error())
			return
		}
	}
	p.connections[conn.GetToken()] = conn
	log.Println("new connection, token:", conn.GetToken())
}

func (p *Pool) LostConnection(conn interfaces.Conn) {
	log.Println("lost connection, token:", conn.GetToken())
	delete(p.connections, conn.GetToken())
	if p.connectedHandler != nil {
		(*p.connectedHandler)(conn)
	}
	conn.Stop()
}

func (p *Pool) HandleLoop() {
	for {
		select {
		case conn := <-p.register:
			p.NewConnection(conn)
		case conn := <-p.unregister:
			if _, ok := p.connections[conn.GetToken()]; ok {
				p.LostConnection(conn)
			}
		}
	}
}