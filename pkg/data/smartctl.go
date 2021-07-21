package data

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Smartctl        Smartctl
	ModelFamily     string      `json:"model_family"`
	ModelName       string      `json:"model_name"`
	SerialNumber    string      `json:"serial_number"`
	FirmwareVersion string      `json:"firmware_version"`
	PowerOnTime     PowerOnTime `json:"power_on_time"`
	PowerCycleCount int         `json:"power_cycle_count"`
	Temperature     Temperature
}

type Smartctl struct {
	Argv       []string
	Output     []string
	ExitStatus int
}

type Device struct {
	Name     string
	InfoName string `json:"info_name"`
	Type     string
	Protocol string
}

type PowerOnTime struct {
	Hours int
}

type Temperature struct {
	Current       int
	PowerCycleMin int `json:"power_cycle_min"`
	PowerCycleMax int `json:"power_cycle_max"`
	LifetimeMin   int `json:"lifetime_min"`
	LifetimeMax   int `json:"lifetime_max"`
	OpLimitMax    int `json:"op_limit_max"`
}

func ParseSmartctlJson(data []byte) (response Response, err error) {
	str := string(data)
	fmt.Print(len(str))

	err = json.Unmarshal(data, &response)
	if err != nil {
		return
	}

	debugJson, _ := json.MarshalIndent(response, "", "    ")
	fmt.Print(string(debugJson))

	fmt.Print("debug")
	return
}
