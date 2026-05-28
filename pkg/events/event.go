package events

type Event interface {
	GetKey() []byte
	GetTopic() string
}
