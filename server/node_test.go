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
	"net"
	"testing"
)

func TestManageNodes(t *testing.T) {
	initServer(&Config{NoBootstrap: true})
	go newNodeManager()
	addr := &net.TCPAddr{net.ParseIP(DefaultIP), DefaultRPCPort, ""}
	addNodeChan <- addr
	addNodeChan <- addr
	if len(nodes) != 1 {
		t.Errorf("after \"addNodeChan <- addr\" twice, len(nodes) = %v", len(nodes))
	}
	removeNodeChan <- addr
	if len(nodes) != 0 {
		t.Errorf("after \"removeNodeChan <- addr\" , len(nodes) = %v", len(nodes))
	}
}
