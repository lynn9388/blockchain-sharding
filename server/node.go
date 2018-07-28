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
	"math/rand"
	"net"
	"net/rpc"
	"sync"

	"github.com/lynn9388/blockchain-sharding/common"
)

var (
	nodes           = make(map[string]common.Node)
	addNodeChan     = make(chan *net.TCPAddr, 5)
	removeNodeChan  = make(chan *net.TCPAddr)
	discoverSigChan = make(chan bool)
	nodeMux         = sync.RWMutex{}

	bootstraps = []net.TCPAddr{
		{net.ParseIP("127.0.0.1"), 9389, ""},
	}
)

func newNodeManager() {
	for _, addr := range bootstraps {
		addNode(&addr)
	}

	for {
		select {
		case addr := <-addNodeChan:
			addNode(addr)
		case addr := <-removeNodeChan:
			removeNode(addr)
		case <-discoverSigChan:
			go discoverNodes()
		}
	}
}

// addNode adds new node to managed node list if it does not exist
func addNode(addr *net.TCPAddr) {
	if !isSelf(addr) {
		nodeMux.Lock()
		if _, exists := nodes[addr.String()]; !exists {
			nodes[addr.String()] = common.Node{*addr}
			common.Logger.Debugf("add new node: %v", addr.String())
		}
		nodeMux.Unlock()
	}
}

// removeNode removes node from managed node list if it exists
func removeNode(addr *net.TCPAddr) {
	nodeMux.Lock()
	if _, exists := nodes[addr.String()]; exists {
		delete(nodes, addr.String())
		common.Logger.Debugf("remove node: %v", addr.String())
	}
	nodeMux.Unlock()
}

func getShuffleNodes() *[]common.Node {
	nodeMux.RLock()
	length := len(nodes)
	tempNodes := make([]common.Node, 0, length)
	for _, node := range nodes {
		tempNodes = append(tempNodes, node)
	}
	nodeMux.RUnlock()

	shuffleNodes := make([]common.Node, length)
	perm := rand.Perm(length)
	for i, v := range perm {
		shuffleNodes[v] = tempNodes[i]
	}
	return &shuffleNodes
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
	for _, n := range *shuffleNodes {
		client, err := connectNode(&n)
		if err != nil {
			removeNodeChan <- &n.RPCAddr
			continue
		}
		newNodes := make([]common.Node, 0)
		err = client.Call("NodeService.GeiNeighborNodes", common.Server.Node.RPCAddr, &newNodes)
		client.Close()
		if err != nil {
			common.Logger.Errorf("failed to call GeiNeighborNodes on %+v: %v", n, err)
		}
		for _, n := range newNodes {
			addNodeChan <- &n.RPCAddr
		}
	}
}
