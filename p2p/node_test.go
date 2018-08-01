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
	"testing"
)

func TestNewNodeManager(t *testing.T) {
	NewNodeManager()
	for _, addr := range bootstraps {
		removeNodeChan <- &addr
	}

	addr := &net.TCPAddr{IP: net.ParseIP("8.8.8.8"), Port: 80}
	addNodeChan <- addr
	if len(getShuffleNodes()) != 1 {
		t.Errorf("after \"addNodeChan <- addr\", len(nodes) = %v", len(getShuffleNodes()))
	}
	addNodeChan <- addr
	if len(getShuffleNodes()) != 1 {
		t.Errorf("after \"addNodeChan <- addr\" twice, len(nodes) = %v", len(getShuffleNodes()))
	}

	removeNodeChan <- addr
	if len(getShuffleNodes()) != 0 {
		t.Errorf("after \"removeNodeChan <- addr\" , len(nodes) = %v", len(getShuffleNodes()))
	}
	removeNodeChan <- addr
	if len(getShuffleNodes()) != 0 {
		t.Errorf("after \"removeNodeChan <- addr\" twice, len(nodes) = %v", len(getShuffleNodes()))
	}
}
