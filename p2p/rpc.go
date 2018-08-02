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
	"net"

	"github.com/lynn9388/blockchain-sharding/common"
	"google.golang.org/grpc"
)

type (
	NodeService int
)

var (
	RPCStartChan = make(chan struct{})
)

func NewRPCListener(rpcAddr string) {
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		common.Logger.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	RegisterDiscoverNodeServer(server, &discoverNodeServer{})

	common.Logger.Infof("start RPC listener on %v", rpcAddr)
	RPCStartChan <- struct{}{}
	server.Serve(lis)
}

func (t *NodeService) GeiNeighborNodes(source *net.TCPAddr, nodes *[]common.Node) error {
	addNodeChan <- &common.Node{RPCAddr: source.String()}
	shuffleNodes := getShuffleNodes()

	length := len(shuffleNodes)
	if maxPeerNum < length {
		length = maxPeerNum
	}
	*nodes = shuffleNodes[:length]
	return nil
}
