/*
 * Copyright 2023 American Express
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

/*
 * Broadcast package is used to create multiple channel reader.
 * Since once we read data from channel, it empties out it self.
 * this package provides a way to broadcast the data and assign listern to that data.
 */
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

// Subscribe() creates a subcribtion on broadcastServer.
func (s *broadcastServer) Subscribe() <-chan scan.Hit {
	newListener := make(chan scan.Hit)
	s.addListener <- newListener
	return newListener
}

// CancelSubscription() cancel a subcribtion on broadcastServer.
func (s *broadcastServer) CancelSubscription(channel <-chan scan.Hit) {
	s.removeListener <- channel
}

// NewBroadcastServer() create a broadcast server and starts new routine.
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

// serve() run the server and manages listener counts.
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
