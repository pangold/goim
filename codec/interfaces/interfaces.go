package interfaces

import (
	message "gitlab.com/pangold/goim/protocol"
)

// It's a message encoder & decoder
// Parse from []byte to segment,
// and callback message
// Use after any bytes received
type Codec interface {
	// Received bytes of data,
	// Pushing into here, and merge
	Encode(interface{}, []byte)
	// Before sending message, split into segments
	Decode(interface{}, *message.Message) error
	// Callback segment, and ready to send
	SetEncodeHandler(func(interface{}, []byte))
	// Callback a received message
	SetDecodeHandler(func(interface{}, *message.Message))
	// Enable resend if timeout
	EnableResend(bool)
}

// It's a received segment pool
// Combine segments to message,
// and callback till a complete message or reply.
// Use after a segment received
type Decoder interface {
	// Received a segment,
	// Pushing into here, and merge
	Push(interface{}, *message.Segment)
	// Callback the complete message
	SetMessageHandler(func(interface{}, *message.Message))
}

// Split into segments if message is too large.
// Use before sending them out,
type Encoder interface {
	// a fake send function,
	// it's real purpose is split into segments if message is too large
	Send(interface{}, *message.Message)
	// use to tell invokers that it's time to send them out
	SetSegmentHandler(func(interface{}, *message.Segment))
	// if it haven't gotten reply for a long time,
	// trigger a resend callback.
	// tips: could be the same to SetPackageHandler
	SetResendHandler(func(interface{}, *message.Segment))
	// every time, invokers received a ack segment,
	// tell me here.
	SetAckSegment(*message.Segment)
}


