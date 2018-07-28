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
	"net/rpc"

	"github.com/lynn9388/blockchain-sharding/common"
)

type (
	PingPongService int
	NodeService     int
)

const (
	pingMsg = "PING"
	pongMsg = "PONG"
)

func newRPCListener(addr *net.TCPAddr) {
	rpc.Register(new(PingPongService))
	rpc.Register(new(NodeService))
	listener, err := net.ListenTCP("tcp", addr)
	defer listener.Close()
	if err != nil {
		common.Logger.Fatalf("failed to start RPC listener: %v", err)
	}
	common.Logger.Infof("start RPC listener on %v", addr.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			common.Logger.Error("failed to accept a RPC connection")
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// PingPong send pong ack message for ping message
func (t *PingPongService) PingPong(msg *string, ack *string) error {
	if *msg != pingMsg {
		return errors.New("not a valid ping message: " + *msg)
	}
	*ack = pongMsg
	return nil
}

func (t *NodeService) GeiNeighborNodes(source *net.TCPAddr, nodes *[]node) error {
	addNodeChan <- source
	shuffleNodes := getShuffleNodes()

	length := len(*shuffleNodes)
	if maxPeerNum < length {
		length = maxPeerNum
	}
	*nodes = (*shuffleNodes)[:length]
	return nil
}
