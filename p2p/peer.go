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

package p2p

import (
	"sync"
	"time"

	"github.com/lynn9388/blockchain-sharding/common"
	"google.golang.org/grpc"
)

const (
	maxPeersNum        = 4
	lackNodesSleepTime = 1
	fullNodesSleepTime = 2
)

type peer struct {
	common.Node
	Conn *grpc.ClientConn
}

var (
	peers    = make(map[string]peer)
	peersMux = sync.RWMutex{}
)

func StartPeerManager() {
	for _, b := range bootstraps {
		addNode(&common.Node{RPCAddr: b})
	}

	go func() {
		for {
			ps := getPeers()
			num := len(ps)
			if num < maxPeersNum {
				sn := getShuffleNodes()
				for i := 0; i < len(sn) && num < maxPeersNum; i++ {
					if _, exists := ps[sn[i].RPCAddr]; exists {
						continue
					}
					if connectNode(sn[i].RPCAddr) {
						num++
					}
				}
			}

			if num < maxPeersNum {
				go discoverNodes()
				time.Sleep(lackNodesSleepTime * time.Second)
			} else {
				time.Sleep(fullNodesSleepTime * time.Second)
			}
		}
	}()
}

// addPeer adds new peer to managed peer list if it does not exist and return true,
func addPeer(p *peer) bool {
	peersMux.Lock()
	defer peersMux.Unlock()
	if _, exists := peers[p.RPCAddr]; !exists {
		peers[p.RPCAddr] = *p
		common.Logger.Debug("added new peer: ", p.RPCAddr)
		return true
	}
	return false
}

// removePeer removes peer from managed peer list if it exists and return true,
func removePeer(rpcAddr string) bool {
	peersMux.Lock()
	defer peersMux.Unlock()
	if _, exists := peers[rpcAddr]; exists {
		peers[rpcAddr].Conn.Close()
		delete(peers, rpcAddr)
		common.Logger.Debug("removed peer: ", rpcAddr)
		return true
	}
	return false
}

// getPeers returns a copy of peers
func getPeers() map[string]peer {
	peersMux.RLock()
	defer peersMux.RUnlock()
	ps := make(map[string]peer)
	for k, v := range peers {
		ps[k] = v
	}
	return ps
}

func connectNode(rpcAddr string) bool {
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		common.Logger.Errorf("failed to dial: %v", err)
		return false
	}
	if !addPeer(&peer{Node: common.Node{RPCAddr: rpcAddr}, Conn: conn}) {
		conn.Close()
		return false
	}
	return true
}
