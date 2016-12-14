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
	"errors"
	"fmt"
	"infra/platformd/objects"
	"infra/platformd/pluginManager/pluginCommon"
	"strconv"
	"time"
	"utils/logging"
)

const (
	NUM_OF_FAN          int = 12
	NUM_OF_THERMAL      int = 8
	SENSOR_POLLING_TIME     = time.Duration(1) * time.Second
)

type openBMCOpen19Driver struct {
	logger logging.LoggerIntf
	ipAddr string
	port   string
}

var driver openBMCOpen19Driver

func NewOpenBMCOpen19Plugin(params *pluginCommon.PluginInitParams) (*openBMCOpen19Driver, error) {
	var err error
	driver.logger = params.Logger
	driver.ipAddr = params.IpAddr
	driver.port = params.Port
	return &driver, err
}

func (driver *openBMCOpen19Driver) Init() error {
	driver.logger.Info("Initializing openBMC driver")
	return nil
}

func (driver *openBMCOpen19Driver) DeInit() error {
	driver.logger.Info("DeInitializing openBMC driver")
	return nil
}

func (driver *openBMCOpen19Driver) GetFanState(fanId int32) (pluginCommon.FanState, error) {
	var state pluginCommon.FanState

	fanDataList, err := driver.getAllFanData()
	if err != nil {
		driver.logger.Err("Error Getting Fan State")
		return state, err
	}
	flag := false
	for _, fanData := range fanDataList {
		if fanId != fanData.FanId {
			continue
		}
		flag = true
		state.Valid = true
		state.FanId = fanData.FanId
		state.OperMode = fanData.OperState
		state.OperSpeed = fanData.OperSpeed
		state.OperDirection = "Not Supported"
		state.Status = fanData.Status
		state.Model = "Not Supported"
		state.SerialNum = "Not Supported"
		state.LedId = -1
	}
	if flag == false {
		return state, errors.New("Unable to find given Fan Id")
	}
	return state, nil
}

func (driver *openBMCOpen19Driver) GetAllFanState(stateList []pluginCommon.FanState, numOfFans int) error {
	fanDataList, err := driver.getAllFanData()
	if err != nil {
		driver.logger.Err("Error Getting Fan State")
		return err
	}
	idx := 0
	for _, fanData := range fanDataList {
		var state pluginCommon.FanState
		state.Valid = true
		state.FanId = fanData.FanId
		state.OperMode = fanData.OperState
		state.OperSpeed = fanData.OperSpeed
		state.OperDirection = "Not Supported"
		state.Status = fanData.Status
		state.Model = "Not Supported"
		state.SerialNum = "Not Supported"
		state.LedId = -1
		stateList[idx] = state
		idx++
	}
	return nil
}

func (driver *openBMCOpen19Driver) GetFanConfig(fanId int32) (retObj *objects.FanConfig, err error) {
	return retObj, nil
}

func (driver *openBMCOpen19Driver) UpdateFanConfig(cfg *objects.FanConfig) (bool, error) {
	driver.logger.Info("Updating OpenBMC Fan Config")
	return true, nil
}

func (driver *openBMCOpen19Driver) GetMaxNumOfFans() int {
	driver.logger.Info("Inside OpenBMC: GetMaxNumOfFans()")
	return NUM_OF_FAN
}

func (driver *openBMCOpen19Driver) GetSfpState(sfpId int32) (pluginCommon.SfpState, error) {
	var retObj pluginCommon.SfpState
	return retObj, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetSfpConfig(sfpId int32) (*objects.SfpConfig, error) {
	var retObj objects.SfpConfig
	return &retObj, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) UpdateSfpConfig(cfg *objects.SfpConfig) (bool, error) {
	driver.logger.Info("Updating  SFP Config")
	return false, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetAllSfpState(states []pluginCommon.SfpState, cnt int) error {
	driver.logger.Info("GetAllSfpState")
	return errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetSfpCnt() int {
	driver.logger.Info("GetSfpCnt")
	return 0
}

func (driver *openBMCOpen19Driver) GetPlatformState() (pluginCommon.PlatformState, error) {
	var retObj pluginCommon.PlatformState
	platformState, err := driver.getPlatformState()
	if err != nil {
		driver.logger.Err("Error Getting Platform State")
		return retObj, err
	}
	retObj.ProductName = platformState.ProductName
	retObj.SerialNum = platformState.ProductVersion
	retObj.Manufacturer = platformState.SystemManufacturer
	retObj.Vendor = platformState.SystemManufacturer
	retObj.Release = fmt.Sprintf("%s.%s", platformState.ProductVersion, platformState.ProductSubVersion)
	//TODO
	retObj.PlatformName = "OpenBMC TBD"
	retObj.Version = "OpenBMC Version TBD"
	return retObj, nil
}

func (driver *openBMCOpen19Driver) GetMaxNumOfThermal() int {
	driver.logger.Info("Inside OpenBMC: GetMaxNumOfThermal()")
	return NUM_OF_THERMAL
}

func (driver *openBMCOpen19Driver) GetThermalState(thermalId int32) (pluginCommon.ThermalState, error) {
	var state pluginCommon.ThermalState
	if thermalId < 0 || thermalId >= int32(NUM_OF_THERMAL) {
		return state, errors.New("Invalid Thermal Id")
	}
	tempDataList, err := driver.getAllTempData()
	if err != nil {
		driver.logger.Err("Error Getting Temperature State")
		return state, err
	}
	flag := false
	for _, tempData := range tempDataList {
		if thermalId != tempData.TempSensorId-1 {
			continue
		}
		flag = true
		state.Valid = true
		state.ThermalId = tempData.TempSensorId - 1
		state.Location = tempData.Name
		state.LowerWatermarkTemperature = "Not Supported"
		state.ShutdownTemperature = "Not Supported"
		state.Temperature = strconv.FormatFloat(tempData.Temperature, 'f', -1, 32) + " C"
		state.UpperWatermarkTemperature = "Not Supported"
	}
	if flag == false {
		state.Valid = false
		state.ThermalId = thermalId
		state.Location = "Unknown"
		state.LowerWatermarkTemperature = "Not Supported"
		state.ShutdownTemperature = "Not Supported"
		state.Temperature = "Unknown"
		state.UpperWatermarkTemperature = "Not Supported"
		return state, nil
	}
	return state, nil
}

func (driver *openBMCOpen19Driver) GetAllThermalState(stateList []pluginCommon.ThermalState, cnt int) error {
	tempDataList, err := driver.getAllTempData()
	if err != nil {
		driver.logger.Err("Error Getting Temperature State")
		return err
	}
	idx := 0
	for _, tempData := range tempDataList {
		var state pluginCommon.ThermalState
		state.Valid = true
		state.ThermalId = tempData.TempSensorId
		state.Location = tempData.Name
		state.LowerWatermarkTemperature = "Not Supported"
		state.ShutdownTemperature = "Not Supported"
		state.Temperature = strconv.FormatFloat(tempData.Temperature, 'f', -1, 32) + " C"
		state.UpperWatermarkTemperature = "Not Supported"
		stateList[idx] = state
		idx++
	}
	return nil
}

func (driver *openBMCOpen19Driver) GetAllSensorState(state *pluginCommon.SensorState) error {
	return errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetQsfpState(Id int32) (retObj pluginCommon.QsfpState, err error) {
	return retObj, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetQsfpPMData(Id int32) (retObj pluginCommon.QsfpPMData, err error) {
	return retObj, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetMaxNumOfQsfp() int {
	driver.logger.Info("GetMaxNumOfQsfps()")
	return 0
}

func (driver *openBMCOpen19Driver) GetPlatformMgmtDeviceState(state *pluginCommon.PlatformMgmtDeviceState) error {
	return errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetMaxNumOfPsu() int {
	return 0
}

func (driver *openBMCOpen19Driver) GetPsuState(psuId int32) (pluginCommon.PsuState, error) {
	var retObj pluginCommon.PsuState
	return retObj, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetAllPsuState(states []pluginCommon.PsuState, cnt int) error {
	return errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetMaxNumOfLed() int {
	return 0
}

func (driver *openBMCOpen19Driver) GetLedState(ledId int32) (pluginCommon.LedState, error) {
	var retObj pluginCommon.LedState
	return retObj, errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetAllLedState(states []pluginCommon.LedState, cnt int) error {
	return errors.New("Not supported")
}

func (driver *openBMCOpen19Driver) GetSfpPortMap() ([]pluginCommon.SfpState, int) {
	return nil, 0
}
