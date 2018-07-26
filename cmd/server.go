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
	"github.com/lynn9388/blockchain-sharding/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A blockchain node and API server",
	Long: `Server is a full blockchain node to connect with other nodes, which will 
construct a peer-to-peer network. It runs a web server to provides REST APIs.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartDaemon(&serverConfig)
	},
}

var serverConfig server.Config

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.
	serverCmd.Flags().StringVarP(&serverConfig.IP, "ip", "i", server.DefaultIP, "the IP address of the server")
	serverCmd.Flags().IntVarP(&serverConfig.APIPort, "api-port", "a", server.DefaultAPIPort, "which port the API service listen on")
	serverCmd.Flags().IntVarP(&serverConfig.RPCPort, "rpc-port", "r", server.DefaultRPCPort, "which port the blockchain node listen on")
	serverCmd.Flags().BoolVar(&serverConfig.NoBootstrap, "no-bootstrap", server.DefaultNoBootstrap, "disable bootstrap nodes on this run")
}
