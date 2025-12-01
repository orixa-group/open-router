package openrouter

// Message is the interface for different message types.
type Message interface {
	Role() string
}
