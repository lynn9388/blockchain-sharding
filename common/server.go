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

package common

import (
	"net"
	"strconv"
)

const (
	DefaultIP      = "127.0.0.1"
	DefaultAPIPort = 9388
	DefaultRPCPort = 9389
)

type (
	Config struct {
		IP      string `json:"ip" description:"ip address of the server" default:"127.0.0.1"`
		APIPort int    `json:"apiport" description:"port of the API Service" default:"9388"`
		RPCPort int    `json:"rpcport" description:"port of the RPC listener" default:"9389"`
	}

	ServerInfo struct {
		APIAddr string
		Node
	}
)

var (
	config     Config
	serverInfo ServerInfo
)

func init() {
	SetConfig(&Config{IP: DefaultIP, APIPort: DefaultAPIPort, RPCPort: DefaultRPCPort})
}

func SetConfig(c *Config) {
	config = *c

	if net.ParseIP(config.IP) == nil {
		Logger.Fatal("failed to parse ip: ", config.IP)
	}
	apiAddr := net.JoinHostPort(config.IP, strconv.Itoa(config.APIPort))
	rpcAddr := net.JoinHostPort(config.IP, strconv.Itoa(config.RPCPort))
	serverInfo = ServerInfo{APIAddr: apiAddr, Node: Node{RPCAddr: rpcAddr}}
}

func GetConfig() *Config {
	return &config
}

func GetServerInfo() *ServerInfo {
	return &serverInfo
}
