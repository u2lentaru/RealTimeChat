package server

type MessageType string
const (
	MTPing    MessageType = "ping"
	MTPong    MessageType = "pong"
	MTMessage MessageType = "message"
)

type Message struct {
	Type MessageType `json:"type"`
	Data string `json:"data,omitempty"`
}
