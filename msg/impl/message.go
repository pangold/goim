package impl

import (
	"github.com/golang/protobuf/proto"
	message "gitlab.com/pangold/goim/msg/protobuf"
	"log"
)

// wrap splitter and merger

type Message struct {
	splitter       *Splitter
	merger         *Merger
	splitHandler   func([]byte)
	messageHandler func(*message.Message)
	ackHandler     func(*message.Message)
	remaining      []byte
}

func NewMessage() *Message {
	p := &Message{
		splitter: NewSplitter(),
		merger:   NewMerger(),
	}
	p.splitter.SetSegmentHandler(p.handleSegment)
	p.splitter.SetResendHandler(p.handleSegment)
	p.merger.SetAckHandler(p.handleAck)
	p.merger.SetMessageHandler(p.handleMessage)
	return p
}

func (p *Message) SetSplitHandler(h func([]byte)) {
	p.splitHandler = h
}

func (p *Message) SetMessageHandler(h func(*message.Message)) {
	p.messageHandler = h
}

func (p *Message) SetAckHandler(h func(*message.Message)) {
	p.ackHandler = h
}

func (p *Message) Split(msg *message.Message) error {
	if err := p.splitter.Send(msg); err != nil {
		return err
	}
	return nil
}

func (p *Message) Merge(data []byte) int {
	seg := &message.Segment{}
	//
	p.remaining = append(p.remaining, data...)
	if err := proto.Unmarshal(p.remaining, seg); err != nil {
		//log.Println(err.Error())
		return seg.XXX_Size()
	}
	// clear remaining
	p.remaining = p.remaining[seg.XXX_Size():]
	// received segment ack
	if seg.GetAck() == 1 {
		if err := p.splitter.SetAckSegment(seg); err != nil {
			log.Println(err.Error())
		}
		return seg.XXX_Size()
	}
	// received a non-ack segment(a real segment)
	if err := p.merger.Push(seg); err != nil {
		log.Println(err.Error())
		return 0
	}
	// reply segment ack
	ack := &message.Segment{
		Id:    proto.Int64(seg.GetId()),
		Index: proto.Int32(0),
		Total: proto.Int32(1),
		Ack:   proto.Int32(1),
		Body:  nil,
	}
	p.handleSegment(ack)
	return seg.XXX_Size()
}

func (p *Message) handleSegment(seg *message.Segment) {
	buf, err := proto.Marshal(seg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	p.splitHandler(buf)
}

func (p *Message) handleResend(seg *message.Segment) {
	p.handleSegment(seg)
}

func (p *Message) handleAck(msg *message.Message) {
	p.handleAck(msg)
}

func (p *Message) handleMessage(msg *message.Message) {
	p.messageHandler(msg)
}
