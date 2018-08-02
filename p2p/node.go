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
	"math/rand"
	"net/rpc"

	"time"

	"github.com/lynn9388/blockchain-sharding/common"
)

var (
	nodes           = make(map[string]common.Node)
	addNodeChan     = make(chan *common.Node)
	removeNodeChan  = make(chan *common.Node)
	getNodesSigChan = make(chan struct{})
	getNodesChan    = make(chan []common.Node)
	discoverSigChan = make(chan bool)

	bootstraps = []string{"127.0.0.1:9388"}
)

func NewNodeManager() {
	for _, addr := range bootstraps {
		addNode(&common.Node{RPCAddr: addr})
	}

	go func() {
		for {
			select {
			case node := <-addNodeChan:
				addNode(node)
			case node := <-removeNodeChan:
				removeNode(node)
			case <-getNodesSigChan:
				getNodesChan <- getNodes()
			case <-discoverSigChan:
				go discoverNodes()
			}
		}
	}()
}

// addNode adds new node to managed node list if it does not exist, it's not safe for concurrent use
func addNode(n *common.Node) {
	if common.Server.RPCAddr != n.RPCAddr {
		if _, exists := nodes[n.RPCAddr]; !exists {
			nodes[n.RPCAddr] = *n
			common.Logger.Debug("add new node: ", n.RPCAddr)
		}
	}
}

// removeNode removes node from managed node list if it exists, it's not safe for concurrent use
func removeNode(n *common.Node) {
	if _, exists := nodes[n.RPCAddr]; exists {
		delete(nodes, n.RPCAddr)
		common.Logger.Debug("remove node: ", n.RPCAddr)
	}
}

// getNodes returns a slice of all nodes, it's not safe for concurrent use
func getNodes() []common.Node {
	n := make([]common.Node, 0, len(nodes))
	for _, node := range nodes {
		n = append(n, node)
	}
	return n
}

// getShuffleNodes returns a shuffled node list, it's safe for concurrent use
func getShuffleNodes() []common.Node {
	getNodesSigChan <- struct{}{}
	nodes := <-getNodesChan

	shuffleNodes := make([]common.Node, len(nodes))
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(nodes))
	for i, v := range perm {
		shuffleNodes[i] = nodes[v]
	}
	return shuffleNodes
}

func connectNode(node *common.Node) (*rpc.Client, error) {
	client, err := rpc.Dial("tcp", node.RPCAddr)
	if err != nil {
		common.Logger.Error(err)
		return nil, err
	}
	return client, nil
}

func discoverNodes() {
	shuffleNodes := getShuffleNodes()
	for _, n := range shuffleNodes {
		client, err := connectNode(&n)
		if err != nil {
			removeNodeChan <- &common.Node{RPCAddr: n.RPCAddr}
			continue
		}
		newNodes := make([]common.Node, 0)
		err = client.Call("NodeService.GeiNeighborNodes", common.Server.RPCAddr, &newNodes)
		client.Close()
		if err != nil {
			common.Logger.Errorf("failed to call GeiNeighborNodes on %+v: %v", n, err)
		}
		for _, n := range newNodes {
			addNodeChan <- &common.Node{RPCAddr: n.RPCAddr}
		}
	}
}
