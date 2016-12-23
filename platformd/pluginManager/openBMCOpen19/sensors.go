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
	"io/ioutil"
	"net/http"
	//"strconv"
	//"strings"
)

type Open19FanState struct {
	FanId     int32  `json:"FanId"`
	OperSpeed int32  `json:"OperationalSpeed"`
	OperState string `json:"OperationalState"`
	Position  string `json:"Position"`
	Status    string `json:"Status"`
}

type Open19FanStateList struct {
	FanStateList []Open19FanState `json:"FanState"`
}

type Open19TempState struct {
	Name         string  `json:"Name"`
	OpeStatus    string  `json:"OperationalStatus"`
	Position     string  `json:"Position"`
	Temperature  float64 `json:"Temperature"`
	TempSensorId int32   `json:"TemperatureSensorID"`
}

type Open19TempStateList struct {
	TempStateList []Open19TempState `json:"TemperatureSensorState"`
}

func SendHttpCmd(req *http.Request) (body []byte, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, err
}

func (driver *openBMCOpen19Driver) getAllFanData() (data []Open19FanState, err error) {
	var jsonStr = []byte(nil)
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/fan"
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
	return extractOpen19FanState(body)
}

func extractOpen19FanState(body []byte) (data []Open19FanState, err error) {
	var info Open19FanStateList

	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("Error:", err)
		return data, err
	}
	return info.FanStateList, err
}

func (driver *openBMCOpen19Driver) getAllTempData() (data []Open19TempState, err error) {
	var jsonStr = []byte(nil)
	url := "http://" + driver.ipAddr + ":" + driver.port + "/api/sys/temp"
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
	return extractOpen19TempState(body)
}

func extractOpen19TempState(body []byte) (data []Open19TempState, err error) {
	var info Open19TempStateList

	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println("Error:", err)
		return data, err
	}
	return info.TempStateList, err
}
