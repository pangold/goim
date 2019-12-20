package tcp

import (
	"gitlab.com/pangold/goim/conn/interfaces"
	"log"
	"net"

	"gitlab.com/im/config"
	"gitlab.com/im/conn"
)

type Server struct {
	conn.ConnectionPool
	receivedHandler     *func([]byte, string) error
	tokenHandler        *func(string) error
	config               config.TcpConfig
	register             chan *Conn
	unregister           chan *Conn
}

func NewTcpServer(c config.TcpConfig) *Server {
	return &Server{
		conn.ConnectionPool {
			nil,
			nil,
			make(map[string]_interface.Connection),
		},
		nil,
		nil,
		c,
		make(chan *Conn),
		make(chan *Conn),
	}
}

func (s *Server) SetReceivedHandler(handler func([]byte, string) error) {
	s.receivedHandler = &handler
}

func (s *Server) SetTokenHandler(handler func(string) error) {
	s.tokenHandler = &handler
}

func (s *Server) Run() {
	//
	ipAddr, err := net.ResolveIPAddr("tcp", s.config.Addr)
	if err != nil {
		log.Fatalf("error: resolve ip address fail: %s\n", s.config.Addr)
		return
	}
	//
	tcpAddr := &net.TCPAddr{
		IP:   ipAddr.IP,
		Port: s.config.Port,
		Zone: ipAddr.Zone,
	}
	//
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("error: listen %s:%d failure. %s\n", s.config.Addr, s.config.Port, err)
		return
	}
	//
	log.Printf("Tcp server start listening %s:%d", s.config.Addr, s.config.Port)
	//
	go s.handleConnection()
	//
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener accept error: ", err)
			continue
		}
		go s.handleAccepted(conn)
	}
}

func (s *Server) handleConnection() {
	for {
		select {
		case conn := <-s.register:
			s.Connections[conn.token] = conn
			if s.ConnectedHandler != nil {
				(*s.ConnectedHandler)(conn.token)
			}
		case conn := <-s.unregister:
			if _, ok := s.Connections[conn.token]; ok {
				close(conn.send)
				delete(s.Connections, conn.token)
				if s.DisconnectedHandler != nil {
					(*s.DisconnectedHandler)(conn.token)
				}
			}
		}
	}
}

func (s *Server) handleAccepted(c net.Conn) {
	conn := &Conn {
		conn:     c,
		received: s.receivedHandler,
		send:     make(chan []byte, 1024),
		token:    "",
	}
	s.register <- conn
	go conn.sendLoop()
	go conn.receiveLoop()
}