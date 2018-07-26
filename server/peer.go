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

import "net/rpc"

func connectNode(node *node) (*rpc.Client, error) {
	client, err := rpc.Dial("tcp", node.rpcAddr.String())
	if err != nil {
		logger.Errorf("failed to dial %v: %v", node.rpcAddr.String(), err)
		return nil, err
	}
	return client, nil
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
