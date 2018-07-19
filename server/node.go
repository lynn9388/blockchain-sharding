/*
 * Copyright © 2018 Lynn <lynn9388@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"errors"
	"net"
	"sync"
)

// A node represents a potential peer on the network
type node struct {
	rpcAddr net.TCPAddr
}

var (
	nodes          = make(map[string]node)
	addNodeChan    = make(chan *net.TCPAddr)
	removeNodeChan = make(chan *net.TCPAddr)
	mux            = sync.RWMutex{}
)

func manageNodes() {
	for {
		select {
		case addr := <-addNodeChan:
			addNode(addr)
		case addr := <-removeNodeChan:
			removeNode(addr)
		}
	}
}

func addNode(addr *net.TCPAddr) error {
	mux.Lock()
	defer mux.Unlock()
	nodes[addr.String()] = node{*addr}
	return nil
}

func removeNode(addr *net.TCPAddr) error {
	mux.Lock()
	defer mux.Unlock()
	if _, exists := nodes[addr.String()]; !exists {
		return errors.New("no record of node: " + addr.String())
	}
	delete(nodes, addr.String())
	return nil
}
