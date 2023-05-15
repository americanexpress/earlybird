package broadcast

import (
	"context"
	"github.com/americanexpress/earlybird/pkg/scan"
)

type BroadcastServer interface {
	Subscribe() <-chan scan.Hit
	CancelSubscription(<-chan scan.Hit)
}

type broadcastServer struct {
	source         <-chan scan.Hit
	listeners      []chan scan.Hit
	addListener    chan chan scan.Hit
	removeListener chan (<-chan scan.Hit)
}

func (s *broadcastServer) Subscribe() <-chan scan.Hit {
	newListener := make(chan scan.Hit)
	s.addListener <- newListener
	return newListener
}

func (s *broadcastServer) CancelSubscription(channel <-chan scan.Hit) {
	s.removeListener <- channel
}

func NewBroadcastServer(ctx context.Context, source <-chan scan.Hit) BroadcastServer {
	service := &broadcastServer{
		source:         source,
		listeners:      make([]chan scan.Hit, 0),
		addListener:    make(chan chan scan.Hit),
		removeListener: make(chan (<-chan scan.Hit)),
	}
	go service.serve(ctx)
	return service
}

func (s *broadcastServer) serve(ctx context.Context) {
	defer func() {
		for _, listener := range s.listeners {
			if listener != nil {
				close(listener)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case newListener := <-s.addListener:
			s.listeners = append(s.listeners, newListener)
		case listenerToRemove := <-s.removeListener:
			for i, ch := range s.listeners {
				if ch == listenerToRemove {
					s.listeners[i] = s.listeners[len(s.listeners)-1]
					s.listeners = s.listeners[:len(s.listeners)-1]
					close(ch)
					break
				}
			}
		case val, ok := <-s.source:
			if !ok {
				return
			}
			for _, listener := range s.listeners {
				if listener != nil {
					select {
					case listener <- val:
					case <-ctx.Done():
						return
					}

				}
			}
		}
	}
}
