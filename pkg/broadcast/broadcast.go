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
	"sync"

	"github.com/americanexpress/earlybird/v4/pkg/scan"
)

type BroadcastServer interface {
	CancelSubscription(<-chan scan.Hit)
	GetListeners() []chan scan.Hit
}

type broadcastServer struct {
	source    <-chan scan.Hit
	listeners []chan scan.Hit
	wg        *sync.WaitGroup
}

// Subscribe() creates a subscription on broadcastServer.
func (s *broadcastServer) Subscribe() <-chan scan.Hit {
	newListener := make(chan scan.Hit)
	s.listeners = append(s.listeners, newListener)

	return newListener
}

// CancelSubscription() cancel a subscription on broadcastServer.
func (s *broadcastServer) CancelSubscription(channel <-chan scan.Hit) {
	for i, ch := range s.listeners {
		if ch == channel {
			s.wg.Done()
			s.listeners[i] = s.listeners[len(s.listeners)-1]
			s.listeners = s.listeners[:len(s.listeners)-1]
			close(ch)
			break
		}
	}
}

func (s *broadcastServer) GetListeners() []chan scan.Hit {
	return s.listeners
}

func (s *broadcastServer) AddSubscriber(count int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		s.wg.Add(1)
		s.Subscribe()
	}
}

func (s *broadcastServer) CloseBroadcast() {
	for _, listener := range s.listeners {
		if listener != nil {
			close(listener)
		}
	}
	// defer s.wg.Done()
}

// NewBroadcastServer() create a broadcast server and starts new routine.
func NewBroadcastServer(ctx context.Context, source <-chan scan.Hit, count int, wg *sync.WaitGroup) BroadcastServer {
	service := &broadcastServer{
		source:    source,
		listeners: make([]chan scan.Hit, 0),
		wg:        wg,
	}
	service.AddSubscriber(count, wg)
	go service.serve(ctx)
	return service
}

// serve() run the server and manages listener counts.
func (s *broadcastServer) serve(ctx context.Context) {
	defer s.CloseBroadcast()

	for {
		select {
		case <-ctx.Done():
			return
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
