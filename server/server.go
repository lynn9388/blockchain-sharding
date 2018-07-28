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
	"os"

	"time"

	"os/signal"
	"syscall"

	"github.com/lynn9388/blockchain-sharding/common"
	"github.com/lynn9388/blockchain-sharding/p2p"
)

func StartServer() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	go newAPIService(&common.Server.APIAddr)
	go p2p.NewRPCListener(&common.Server.Node.RPCAddr)

	time.Sleep(2 * time.Second)

	go p2p.NewNodeManager()
	go p2p.NewPeerManager()

	select {
	case <-sigChan:
		common.Logger.Info("caught stop signal, quitting...")
	}
}
