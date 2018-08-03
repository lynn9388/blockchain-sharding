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
	"context"
	"io"
	"net"
	"strconv"
	"testing"

	"github.com/lynn9388/blockchain-sharding/common"
	"google.golang.org/grpc"
)

func createDiscoverNodeServer(t *testing.T) (*net.Listener, *grpc.Server) {
	lis, err := net.Listen("tcp", common.GetServerInfo().RPCAddr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	RegisterDiscoverNodeServer(server, &discoverNodeServer{})

	return &lis, server
}

func TestDiscoverNodeServer_Ping(t *testing.T) {
	lis, server := createDiscoverNodeServer(t)
	defer server.Stop()

	go server.Serve(*lis)

	conn, err := grpc.Dial(common.GetServerInfo().RPCAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := NewDiscoverNodeClient(conn)
	pong, err := client.Ping(context.Background(), &PingPong{Message: PingPong_PING})
	if err != nil {
		t.Fatal(err)
	}
	if pong.Message != PingPong_PONG {
		t.Fatalf("invalid pong message: %v", pong.Message)
	}
}

func TestDiscoverNodeServer_GeiNeighborNodes(t *testing.T) {
	for port := 80; port < 80+maxShareNodesNum; port++ {
		addNode(&common.Node{RPCAddr: net.JoinHostPort("8.8.8.8", strconv.Itoa(port))})
	}

	lis, server := createDiscoverNodeServer(t)
	defer server.Stop()

	go server.Serve(*lis)

	conn, err := grpc.Dial(common.GetServerInfo().RPCAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := NewDiscoverNodeClient(conn)
	stream, err := client.GeiNeighborNodes(context.Background(), &common.GetServerInfo().Node)
	if err != nil {
		t.Fatal(err)
	}

	ns := make(map[string]*common.Node)
	for {
		n, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		ns[n.RPCAddr] = n
	}
	if len(ns) != maxShareNodesNum {
		t.FailNow()
	}
}
