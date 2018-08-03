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
	"context"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/lynn9388/blockchain-sharding/common"
)

var (
	nodes    = make(map[string]common.Node)
	nodesMux = sync.RWMutex{}

	bootstraps = []string{"127.0.0.1:9389"}
)

// addNode adds new node to managed node list if it does not exist
func addNode(n *common.Node) {
	nodesMux.Lock()
	defer nodesMux.Unlock()
	if common.GetServerInfo().RPCAddr != n.RPCAddr {
		if _, exists := nodes[n.RPCAddr]; !exists {
			nodes[n.RPCAddr] = *n
			common.Logger.Debug("added new node: ", n.RPCAddr)
		}
	}
}

// removeNode removes node from managed node list if it exists
func removeNode(n *common.Node) {
	nodesMux.Lock()
	defer nodesMux.Unlock()
	if _, exists := nodes[n.RPCAddr]; exists {
		delete(nodes, n.RPCAddr)
		common.Logger.Debug("removed node: ", n.RPCAddr)
	}
}

// getNodes returns a copy of nodes
func getNodes() map[string]common.Node {
	nodesMux.RLock()
	defer nodesMux.RUnlock()
	ns := make(map[string]common.Node)
	for k, v := range nodes {
		ns[k] = v
	}
	return ns
}

// getShuffleNodes returns a shuffled node list
func getShuffleNodes() []common.Node {
	ns := getNodes()

	shuffleNodes := make([]common.Node, len(ns))
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(ns))

	i := 0
	for _, v := range ns {
		shuffleNodes[perm[i]] = v
		i++
	}

	return shuffleNodes
}

func isBootstrap(addr string) bool {
	for _, b := range bootstraps {
		if addr == b {
			return true
		}
	}
	return false
}

func discoverNodes() {
	ns := getNodes()
	ps := getPeers()
	for k, v := range ps {
		if _, exists := ns[k]; !exists {
			ns[k] = common.Node{RPCAddr: v.RPCAddr}
			addNode(&common.Node{RPCAddr: v.RPCAddr})
		}
	}

	for _, p := range ps {
		client := NewDiscoverNodeClient(p.Conn)
		stream, err := client.GeiNeighborNodes(context.Background(), &common.GetServerInfo().Node)
		if err != nil {
			common.Logger.Error(err)
			removeNode(&p.Node)
			removePeer(p.RPCAddr)
			continue
		}

		for {
			n, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				common.Logger.Error(err)
				removeNode(&p.Node)
				removePeer(p.RPCAddr)
				break
			}
			if _, exists := ns[n.RPCAddr]; !exists {
				ns[n.RPCAddr] = *n
				addNode(n)
			}
		}
	}

	if len(ns) == 0 {
		for _, b := range bootstraps {
			addNode(&common.Node{RPCAddr: b})
		}
	}
}
