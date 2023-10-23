package helpers_events

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

var lock = &sync.Mutex{}
var dispatcherInstance *Dispatcher

func GetDispatcher() *Dispatcher {
	if dispatcherInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if dispatcherInstance == nil {
			d := &Dispatcher{
				jobs:      make(chan job),
				listeners: make(map[string][]ListenerHandleFunc),
			}
			go d.consume()
			dispatcherInstance = d
		}
	}
	return dispatcherInstance
}

type Dispatcher struct {
	jobs      chan job
	listeners map[string][]ListenerHandleFunc
}

func (x *Dispatcher) Register(name string, listeners ...ListenerHandleFunc) error {
	for _, listener := range listeners {
		x.listeners[name] = append(x.listeners[name], listener)
	}
	return nil
}
func (x *Dispatcher) Dispatch(name string, eventData any) error {
	_, ok := x.listeners[name]
	if !ok {
		return fmt.Errorf("event '%s' is not registered", name)
	}
	logrus.WithField("event", name).Info()
	x.jobs <- job{
		eventName: name,
		eventData: eventData,
	}
	return nil
}
func (x *Dispatcher) consume() {
	for j := range x.jobs {
		for _, listener := range x.listeners[j.eventName] {
			listener(j.eventData)
		}
	}
}
