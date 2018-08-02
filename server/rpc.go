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

	"github.com/lynn9388/blockchain-sharding/common"
	"github.com/lynn9388/blockchain-sharding/p2p"
	"google.golang.org/grpc"
)

var rpcServer *grpc.Server

func startRPCServer() {
	lis, err := net.Listen("tcp", common.GetServerInfo().RPCAddr)
	if err != nil {
		common.Logger.Fatalf("failed to listen: %v", err)
	}

	rpcServer = grpc.NewServer()
	p2p.RegisterDiscoverNodeServer(rpcServer, p2p.NewDiscoverNodeServer())

	common.Logger.Infof("RPC server listening at: %v", common.GetServerInfo().RPCAddr)
	go rpcServer.Serve(lis)
}

func stopRPCServer() {
	if rpcServer != nil {
		rpcServer.Stop()
		common.Logger.Info("RPC server stopped")
	}
}
