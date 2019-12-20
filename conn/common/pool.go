package common

import (
	"gitlab.com/pangold/goim/conn/interfaces"
	"log"
)

type Pool struct {
	connections          map[string]interfaces.Connection
	connectedHandler    *func(interfaces.Connection) error
	disconnectedHandler *func(interfaces.Connection) error
	register             chan interfaces.Connection
	unregister           chan interfaces.Connection
}

func NewPool() interfaces.Pool {
	return &Pool{
		connections:         make(map[string]interfaces.Connection),
		connectedHandler:    nil,
		disconnectedHandler: nil,
		register:            make(chan interfaces.Connection),
		unregister:          make(chan interfaces.Connection),
	}
}

// for upstream invokers
func (p *Pool) SetConnectedHandler(fn func(interfaces.Connection) error) {
	p.connectedHandler = &fn
}

func (p *Pool) SetDisconnectedHandler(fn func(interfaces.Connection) error) {
	p.disconnectedHandler = &fn
}

func (p *Pool) GetConnections() *map[string]interfaces.Connection {
	return &p.connections
}

// for the same level invokers
func (p *Pool) IsExist(conn interfaces.Connection) bool {
	if _, ok := p.connections[conn.GetToken()]; ok {
		return true
	}
	return false
}

func (p *Pool) Register(conn interfaces.Connection) {
	p.register <- conn
}

func (p *Pool) Unregister(conn interfaces.Connection) {
	p.unregister <- conn
}

func (p *Pool) NewConnection(conn interfaces.Connection) {
	if p.connectedHandler != nil {
		// check if it is valid connection(actually, check for token)
		if err := (*p.connectedHandler)(conn); err != nil {
			log.Fatalln(err.Error())
		} else {
			p.connections[conn.GetToken()] = conn
			log.Println("new connection, count: ", len(p.connections))
		}
	} else {
		p.connections[conn.GetToken()] = conn
		log.Println("new connection, count: ", len(p.connections))
	}
}

func (p *Pool) LostConnection(conn interfaces.Connection) {
	log.Println("lost connection, token: ", conn.GetToken())
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