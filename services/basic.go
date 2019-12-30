package services

import (
	"gitlab.com/pangold/goim/codec"
	"gitlab.com/pangold/goim/config"
	"gitlab.com/pangold/goim/conn"
	"gitlab.com/pangold/goim/session"
	"log"
)

type common struct {
	conn     *conn.Server
	codec    *codec.Codec
	sessions *session.SessionPool
}

func NewCommon(conf config.Config) *common {
	c := &common{
		conn:     conn.NewServer(conf),
		codec:    codec.NewCodec(),
		sessions: session.NewSessionPool(),
	}
	c.codec.SetEncodeHandler(c.handleEncode)
	c.codec.SetDecodeHandler(c.handleDecode)
	c.conn.SetMessageHandler(c.handleReceived)
	c.conn.SetConnectedHandler(c.handleConnected)
	c.conn.SetDisconnectedHandler(c.handleDisconnected)
	return c
}

func (c *common) handleConnected(token string) error {
	if err := c.sessions.Push(token); err != nil {
		return err
	}
	return nil
}

func (c *common) handleDisconnected(token string) {
	c.sessions.Remove(token)
}

func (c *common) handleReceived(data []byte, token string) error {
	c.codec.Decode(token, data) // returned value is unused.
	return nil
}

func (c *common) handleDecode(token interface{}, msg *codec.MessageT) {
	// your business model.
}

func (c *common) handleEncode(token interface{}, data []byte) {
	if err := c.conn.Send(token.(string), data); err != nil {
		log.Printf(err.Error())
	}
}

func (c *common) send(token string, msg *codec.MessageT) {
	if err := c.codec.Encode(token, msg); err != nil {
		log.Printf(err.Error())
	}
}
