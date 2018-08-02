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
	"errors"

	"github.com/lynn9388/blockchain-sharding/common"
	"golang.org/x/net/context"
)

const (
	maxShareNodesNum = 10
)

type discoverNodeServer struct {
}

func (s *discoverNodeServer) Ping(ctx context.Context, ping *PingPong) (*PingPong, error) {
	if ping.Message != PingPong_PING {
		return nil, errors.New("invalid ping message: " + ping.Message.String())
	}
	return &PingPong{Message: PingPong_PONG}, nil
}

func (s *discoverNodeServer) GeiNeighborNodes(n *common.Node, stream DiscoverNode_GeiNeighborNodesServer) error {
	addNode(n)

	ns := getShuffleNodes()
	count := 0
	for _, n := range ns {
		if count >= maxShareNodesNum {
			break
		}
		if !isBootstrap(n.RPCAddr) {
			if err := stream.Send(&n); err != nil {
				return err
			}
			count++
		}
	}
	return nil
}
