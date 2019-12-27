package interfaces

import (
	message "gitlab.com/pangold/goim/msg/protobuf"
)

// It's a received buffer
// Parse from []byte to segment,
// and callback
// Use after any bytes received
type Buffer interface {
	// Error info
	Error() error
	// Received bytes data,
	// Pushing into here, and recording
	Push([]byte)
	// Callback the segment
	SetSegmentHandler(func(*message.Segment))
	// Callback the ack
	SetAckHandler(func(*message.Segment))
}

// It's a received segment pool
// Combine segments to message,
// and callback till a complete message or reply.
// Use after a segment received
type Combiner interface {
	// Error information
	Error() error
	// Received a segment,
	// Pushing into here, and recording
	Push(*message.Segment)
	// Callback the complete message
	// Either body or replied field is valid
	// Body: means it is a message from other connection
	// Replied: means it is replied, could be operated state(read, download...)
	SetMessageHandler(func(m *message.Message))
	// Message acknowledge(
	SetAckHandler(func(m *message.Message))
}

// Split into segments if message is too large.
// Use before sending them out,
type Splitor interface {
	// get error information
	Error() error
	// a fake send function,
	// it's real purpose is split into segments if message is too large
	Send(*message.Message)
	// use to tell invokers that it's time to send them out
	SetSegmentHandler(*message.Segment)
	// if it haven't gotten reply for a long time,
	// trigger a resend callback.
	// tips: could be the same to SetPackageHandler
	SetResendHandler(func(*message.Segment))
	// every time, invokers received a ack segment,
	// tell me here.
	SetAckSegment(*message.Segment)
}


