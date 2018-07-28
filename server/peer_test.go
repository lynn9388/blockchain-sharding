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
	"time"

	"github.com/lynn9388/blockchain-sharding/common"
)

func TestPing(t *testing.T) {
	addr := net.TCPAddr{IP: net.ParseIP(common.DefaultIP), Port: common.DefaultRPCPort}
	go newRPCListener(&addr)
	time.Sleep(1 * time.Second)
	if client := ping(&common.Node{RPCAddr: addr}); client == nil {
		t.FailNow()
	}
}
