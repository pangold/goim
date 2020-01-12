package tcp

import (
	"log"
	"net"

	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/front/interfaces"
)

type Server struct {
	config               config.HostConfig
	pool                 interfaces.Pool
}

func NewTcpServer(p interfaces.Pool, c config.HostConfig) *Server {
	return &Server{
		config: c,
		pool:   p,
	}
}

func (s *Server) Run() {
	go s.pool.HandleLoop()
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.config.Address)
	if err != nil {
		log.Fatalf("error: failed to resolve ip address: tcp://%s", s.config.Address)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("error: failed to listen tcp://%s, %v", s.config.Address, err)
		return
	}
	log.Printf("Tcp server start listening tcp://%s", s.config.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener accept error: ", err)
			continue
		}
		go s.handleAccepted(conn)
	}
}

func (s *Server) handleAccepted(c net.Conn) {
	conn := &Connection{
		pool:           s.pool,
		conn:           c,
		messageHandler: nil,
		send:           make(chan []byte, 1024),
		token:          "",
		stopped:        false,
		remaining:      nil,
	}
	// s.pool.Register(conn) // when connection received token
	go conn.sendLoop()
	go conn.receiveLoop()
}