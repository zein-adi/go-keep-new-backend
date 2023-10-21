package helpers_events

type job struct {
	eventName string
	eventData any
}
type ListenerHandleFunc func(eventData any)
