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
	"testing"

	"net"
	"strconv"

	"github.com/lynn9388/blockchain-sharding/common"
)

func TestAddNode(t *testing.T) {
	nodes = make(map[string]common.Node)

	node := &common.Node{RPCAddr: "8.8.8.8:80"}
	addNode(node)
	if len(nodes) != 1 {
		t.Errorf("after \"addNode(node)\", len(nodes) = %v", len(nodes))
	}
	addNode(node)
	if len(nodes) != 1 {
		t.Errorf("after \"addNode(node)\" twice, len(nodes) = %v", len(nodes))
	}
}

func TestRemoveNode(t *testing.T) {
	nodes = make(map[string]common.Node)

	node := &common.Node{RPCAddr: "8.8.8.8:80"}
	addNode(node)
	removeNode(node)
	if len(nodes) != 0 {
		t.Errorf("after \"removeNode(node)\" , len(nodes) = %v", len(nodes))
	}
	removeNode(node)
	if len(nodes) != 0 {
		t.Errorf("after \"removeNode(node)\" twice, len(nodes) = %v", len(nodes))
	}
}

func TestGetNodes(t *testing.T) {
	nodes = make(map[string]common.Node)

	if len(getNodes()) != 0 {
		t.Errorf("\"getNodes()\" from empty nodes, len(getNodes()) = %v", len(getNodes()))
	}

	node := &common.Node{RPCAddr: "8.8.8.8:80"}
	addNode(node)
	if len(getNodes()) != 1 {
		t.Errorf("\"getNodes()\" from nodes (1 node), len(getNodes()) = %v", len(getNodes()))
	}
}

func TestGetShuffleNodes(t *testing.T) {
	nodes = make(map[string]common.Node)

	num := 10
	for port := 80; port < 80+num; port++ {
		addNode(&common.Node{RPCAddr: net.JoinHostPort("8.8.8.8", strconv.Itoa(port))})
	}
	sns1 := getShuffleNodes()
	sns2 := getShuffleNodes()

	diff := 0
	for i := 0; i < num; i++ {
		if sns1[i].RPCAddr != sns2[i].RPCAddr {
			diff++
		}
	}
	if float32(diff)/float32(num) < 0.5 {
		t.FailNow()
	}
}
