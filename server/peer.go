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
	"net/rpc"
	"sync"
)

type peer struct {
	node   *node
	client *rpc.Client
}

const (
	maxPeerNum = 4
)

var (
	peers          = make(map[string]peer)
	addPeerChan    = make(chan *node)
	removePeerChan = make(chan *node)
	peerMux        = sync.RWMutex{}
)

func newPeerManager() {
	for {
		select {
		case node := <-addPeerChan:
			addPeer(node)
		case node := <-removePeerChan:
			removePeer(node)
		}
	}
}

func addPeer(node *node) {
	client := ping(node)
	if client != nil {
		peerMux.Lock()
		if _, exists := peers[node.rpcAddr.String()]; exists {
			peers[node.rpcAddr.String()].client.Close()
		}
		peers[node.rpcAddr.String()] = peer{node, client}
		peerMux.Unlock()
	}
}

func removePeer(node *node) {
	peerMux.Lock()
	peers[node.rpcAddr.String()].client.Close()
	delete(peers, node.rpcAddr.String())
	peerMux.Unlock()
}

// ping tests if a node is reachable and returns connected client
func ping(node *node) *rpc.Client {
	ack := ""
	client, err := connectNode(node)
	if err != nil {
		return nil
	}
	err = client.Call("PingPongService.PingPong", pingMsg, &ack)
	if err != nil {
		logger.Errorf("failed to call PingPong on %v: %v", node.rpcAddr.String(), err)
		return nil
	}
	if ack != pongMsg {
		logger.Errorf("not a valid pong message: %v", ack)
		return nil
	}
	return client
}
