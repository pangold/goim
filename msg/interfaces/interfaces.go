package interfaces

import (
	message "gitlab.com/pangold/goim/msg/protobuf"
)

// It's a message
// Parse from []byte to segment,
// and callback
// Use after any bytes received
type Message interface {
	// Received bytes of data,
	// Pushing into here, and merge
	Merge([]byte) int
	// Before sending message, split into segments
	Split(msg *message.Message) error
	// Callback segment, and ready to send
	SetSplitHandler(func([]byte))
	// Callback a received message
	SetMessageHandler(func(*message.Message))
	// Callback a ack message
	SetAckHandler(func(*message.Message))
}

// It's a received segment pool
// Combine segments to message,
// and callback till a complete message or reply.
// Use after a segment received
type Merger interface {
	// Received a segment,
	// Pushing into here, and merge
	Push(*message.Segment)
	// Callback the complete message
	SetMessageHandler(func(m *message.Message))
	// Message acknowledge
	SetAckHandler(func(m *message.Message))
}

// Split into segments if message is too large.
// Use before sending them out,
type Splitter interface {
	// a fake send function,
	// it's real purpose is split into segments if message is too large
	Send(*message.Message)
	// use to tell invokers that it's time to send them out
	SetSegmentHandler(func(*message.Segment))
	// if it haven't gotten reply for a long time,
	// trigger a resend callback.
	// tips: could be the same to SetPackageHandler
	SetResendHandler(func(*message.Segment))
	// every time, invokers received a ack segment,
	// tell me here.
	SetAckSegment(*message.Segment)
}


