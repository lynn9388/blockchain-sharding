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
	"net"
	"net/rpc"

	"time"

	"github.com/lynn9388/blockchain-sharding/common"
)

var (
	nodes           = make(map[string]common.Node)
	addNodeChan     = make(chan *net.TCPAddr)
	removeNodeChan  = make(chan *net.TCPAddr)
	getNodesSigChan = make(chan struct{})
	getNodesChan    = make(chan []common.Node)
	discoverSigChan = make(chan bool)

	bootstraps = []net.TCPAddr{
		{IP: net.ParseIP("127.0.0.1"), Port: 9389},
	}
)

func NewNodeManager() {
	for _, addr := range bootstraps {
		addNode(&addr)
	}

	for {
		select {
		case addr := <-addNodeChan:
			addNode(addr)
		case addr := <-removeNodeChan:
			removeNode(addr)
		case <-getNodesSigChan:
			getNodesChan <- getNodes()
		case <-discoverSigChan:
			go discoverNodes()
		}
	}
}

// addNode adds new node to managed node list if it does not exist, it's not safe for concurrent use
func addNode(addr *net.TCPAddr) {
	if common.Server.RPCAddr.String() != addr.String() {
		if _, exists := nodes[addr.String()]; !exists {
			nodes[addr.String()] = common.Node{RPCAddr: *addr}
			common.Logger.Debug("add new node: ", addr.String())
		}
	}
}

// removeNode removes node from managed node list if it exists, it's not safe for concurrent use
func removeNode(addr *net.TCPAddr) {
	if _, exists := nodes[addr.String()]; exists {
		delete(nodes, addr.String())
		common.Logger.Debug("remove node: ", addr.String())
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
	client, err := rpc.Dial("tcp", node.RPCAddr.String())
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
			removeNodeChan <- &n.RPCAddr
			continue
		}
		newNodes := make([]common.Node, 0)
		err = client.Call("NodeService.GeiNeighborNodes", common.Server.RPCAddr, &newNodes)
		client.Close()
		if err != nil {
			common.Logger.Errorf("failed to call GeiNeighborNodes on %+v: %v", n, err)
		}
		for _, n := range newNodes {
			addNodeChan <- &n.RPCAddr
		}
	}
}
