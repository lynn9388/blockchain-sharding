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

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	DefaultIP      = "127.0.0.1"
	DefaultAPIPort = 9388
	DefaultRPCPort = 9389
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A blockchain node and API server",
	Long: `Server is a full blockchain node to connect with other nodes, which will 
construct a peer-to-peer network. It runs a web server to provides REST APIs.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ServerLogger.Debugf("server's configuration: %+v", ServerConfig)
	},
}

type Server struct {
	IP      string `json:"ip" description:"ip address of the server" default:"127.0.0.1"`
	APIPort int    `json:"apiport" description:"port of the API Service" default:"9388"`
	RPCPort int    `json:"rpcport" description:"port of the RPC listener" default:"9389"`
}

var (
	ServerConfig Server
	ServerLogger *zap.SugaredLogger
)

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.
	serverCmd.Flags().StringVarP(&ServerConfig.IP, "ip", "i", DefaultIP, "the IP address of the server")
	serverCmd.Flags().IntVarP(&ServerConfig.APIPort, "api-port", "a", DefaultAPIPort, "which port the API server listen on")
	serverCmd.Flags().IntVarP(&ServerConfig.RPCPort, "rpc-port", "r", DefaultRPCPort, "which port the blockchain node listen on")

	logger, _ := zap.NewDevelopment()
	ServerLogger = logger.Sugar()
}
