package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn/interfaces"
	"log"
	"net/http"
)

type Server struct {
	upgrader             websocket.Upgrader
	config               config.WsConfig
	pool                 interfaces.Pool
}

func UseUpgrader(r, w int) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize: r,
		WriteBufferSize: w,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
}

func NewWsServer(p interfaces.Pool, c config.WsConfig) *Server {
	return &Server{
		upgrader:        UseUpgrader(1024, 1024),
		config:          c,
		pool:            p,
	}
}

func (s *Server) Run() {
	go s.pool.HandleLoop()
	//
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ws", s.handleWs)
	log.Printf("Websocket server start running with socket protocol: %s, serving: %s", s.config.Protocol, s.config.Address)
	//
	if s.config.Protocol == "wss" {
		if err := http.ListenAndServeTLS(s.config.Address, s.config.CertFile, s.config.KeyFile, serverMux); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(s.config.Address, serverMux); err != nil {
			panic(err)
		}
	}
}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	token, err := s.checkToken(r.URL.Query())
	if err != nil {
		log.Fatal(err.Error())
		// Unauthorized
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websocket upgrade error: %s", err.Error())
		return
	}
	conn := &Connection{
		pool:           s.pool,
		conn:           c,
		messageHandler: nil,
		send:           make(chan []byte, 1024),
		token:          token,
	}
	if s.pool.IsExist(conn) {
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	s.pool.Register(conn)
	go conn.sendLoop()
	go conn.receiveLoop()
}

func (s Server) checkToken(query map[string][]string) (string, error) {
	if tokens, ok := query["token"]; ok && len(tokens) > 0 && len(tokens[0]) > 0 {
		return tokens[0], nil
	}
	return "", errors.New("error: unmatched token param")
}
