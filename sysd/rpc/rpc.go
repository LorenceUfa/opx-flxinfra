//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package rpc

import (
	"encoding/json"
	"git.apache.org/thrift.git/lib/go/thrift"
	"io/ioutil"
	"strconv"
	"sysd"
	"utils/logging"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

func getClient(logger *logging.Writer, fileName string, process string) (*ClientJson, error) {
	var allClients []ClientJson

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Err("Failed to open config file:%s, err:%s", fileName, err)
		return nil, err
	}

	json.Unmarshal(data, &allClients)
	for _, client := range allClients {
		if client.Name == process {
			return &client, nil
		}
	}

	logger.Err("Did not find port for %s in config file:%s", process, fileName)
	return nil, nil
}

func StartServer(logger *logging.Writer, handler *SYSDHandler, fileName string) {
	clientJson, err := getClient(logger, fileName, "sysd")
	if err != nil || clientJson == nil {
		return
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	serverTransport, err := thrift.NewTServerSocket("localhost:" + strconv.Itoa(clientJson.Port))
	if err != nil {
		logger.Info("StartServer: NewTServerSocket failed with error:", err)
		return
	}
	processor := sysd.NewSYSDServicesProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	logger.Info("StartServer: Starting")
	err = server.Serve()
	logger.Info("StartServer: Started")
	if err != nil {
		logger.Info("Failed to start the listener, err:", err)
	}
	logger.Info("Started the listener successfully")
	return
}