package server

// MessageType type
type MessageType string

// MTPing, MTPong, MTMessage const
const (
	MTPing    MessageType = "ping"
	MTPong    MessageType = "pong"
	MTMessage MessageType = "message"
)

// Message type
type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data,omitempty"`
}
