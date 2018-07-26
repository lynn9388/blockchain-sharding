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
	"strconv"

	"os"

	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type (
	Config struct {
		IP          string `json:"ip" description:"ip address of the server" default:"127.0.0.1"`
		APIPort     int    `json:"apiport" description:"port of the API Service" default:"9388"`
		RPCPort     int    `json:"rpcport" description:"port of the RPC listener" default:"9389"`
		NoBootstrap bool   `json:"no-bootstrap" description:"disable bootstrap nodes on this run" default:"false"`
	}

	server struct {
		apiAddr net.TCPAddr
		node    node
	}
)

const (
	DefaultIP          = "127.0.0.1"
	DefaultAPIPort     = 9388
	DefaultRPCPort     = 9389
	DefaultNoBootstrap = false
)

var (
	logger  *zap.SugaredLogger
	sigChan chan os.Signal
	config  *Config
	daemon  server
)

func init() {
	l, _ := zap.NewDevelopment()
	logger = l.Sugar()

	sigChan = make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
}

func StartDaemon(c *Config) {
	initServer(config)

	go newAPIService(&daemon.apiAddr)
	go newRPCListener(&daemon.node.rpcAddr)
	go newNodeManager()

	select {
	case <-sigChan:
		logger.Info("caught stop signal, quitting...")
	}
}

// initServer initials config and combines API service address and RPC listener address from configuration
func initServer(c *Config) {
	config = c

	ip := config.IP
	apiPort := strconv.Itoa(config.APIPort)
	rpcPort := strconv.Itoa(config.RPCPort)

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(ip, apiPort))
	if err != nil {
		logger.Fatalf("failed to combine API service address: %+v", err)
	}
	daemon.apiAddr = *addr

	addr, err = net.ResolveTCPAddr("tcp", net.JoinHostPort(ip, rpcPort))
	if err != nil {
		logger.Fatalf("failed to combine RPC listener address: %+v", err)
	}
	daemon.node = node{*addr}
}

func isSelf(addr *net.TCPAddr) bool {
	return addr.String() == daemon.node.rpcAddr.String()
}
