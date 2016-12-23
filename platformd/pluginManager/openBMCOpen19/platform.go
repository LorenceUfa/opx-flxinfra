//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package openBMCOpen19

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	//"strconv"
	//"strings"
)

type Open19PlatformStateData struct {
	AssembledAt              string `json:"Assembled At"`
	CRC8                     string `json:"CRC8"`
	ExtendedMACAddrSize      string `json:"Extended MAC Address Size"`
	ExtendedMACBase          string `json:"Extended MAC Base"`
	FlexPCBPartNumber        string `json:"Flex PCB Part Number"`
	FlexPCBAPartNumber       string `json:"Flex PCBA Part Number"`
	HawkEEPROM               string `json:"Hawk EEPROM"`
	LocalMAC                 string `json:"Local MAC"`
	LocationOnFabric         string `json:"Location on Fabric"`
	ODMPCBAPartNumber        string `json:"ODM PCBA Part Number"`
	ODMPCBASerialNumber      string `json:"ODM PCBA Serial Number"`
	PCBManufacturer          string `json:"PCB Manufacturer"`
	ProductAssetTag          string `json:"Product Asset Tag"`
	ProductName              string `json:"Product Name"`
	ProductPartNumber        string `json:"Product Part Number"`
	ProductProductionState   string `json:"Product Production State"`
	ProductSerialNumber      string `json:"Product Serial Number"`
	ProductSubVersion        string `json:"Product Sub-Version"`
	ProductVersion           string `json:"Product Version"`
	SystemAssemblyPartNumber string `json:"System Assembly Part Number"`
	SystemManufacturer       string `json:"System Manufacturer"`
	SystemManufacturingDate  string `json:"System Manufacturing Date"`
	Version                  string `json:"Version"`
}

type Open19PlatformState struct {
	PlatformState []Open19PlatformStateData `json:"PlatformState"`
}

func (driver *openBMCOpen19Driver) getPlatformState() (data Open19PlatformStateData, err error) {
	var jsonStr = []byte(nil)
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/platform"
	//fmt.Println("URL:>", url)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return data, err
	}
	req.Header.Set("Accept", "application/json")
	body, err := SendHttpCmd(req)
	if err != nil {
		return data, err
	}
	return extractOpen19PlatformState(body)
}

func extractOpen19PlatformState(body []byte) (data Open19PlatformStateData, err error) {
	var info Open19PlatformState

	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("Error:", err)
		return data, err
	}
	return info.PlatformState[0], err
}
