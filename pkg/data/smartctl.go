package data

import (
	"encoding/json"
)

type Response struct {
	Smartctl           Smartctl           `json:"smartctl"`
	ModelFamily        string             `json:"model_family"`
	ModelName          string             `json:"model_name"`
	SerialNumber       string             `json:"serial_number"`
	FirmwareVersion    string             `json:"firmware_version"`
	PowerOnTime        PowerOnTime        `json:"power_on_time"`
	PowerCycleCount    int                `json:"power_cycle_count"`
	Temperature        Temperature        `json:"temperature"`
	AtaSmartSttributes AtaSmartSttributes `json:"ata_smart_attributes"`
	AtaSctStatus       AtaSctStatus       `json:"ata_sct_status"`
	SmartStatus        SmartStatus        `json:"smart_status"`
}

type Smartctl struct {
	Argv       []string `json:"argv"`
	Output     []string `json:"output"`
	ExitStatus int      `json:"exit_status"`
}

type Device struct {
	Name     string `json:"name"`
	InfoName string `json:"info_name"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

type SmartStatus struct {
	Passed bool `json:"passed"`
}

type PowerOnTime struct {
	Hours int `json:"hours"`
}

type Temperature struct {
	Current       int `json:"current"`
	PowerCycleMin int `json:"power_cycle_min"`
	PowerCycleMax int `json:"power_cycle_max"`
	LifetimeMin   int `json:"lifetime_min"`
	LifetimeMax   int `json:"lifetime_max"`
	OpLimitMax    int `json:"op_limit_max"`
}

type AtaSmartSttributes struct {
	Revision int              `json:"revision"`
	Table    []SmartAttribute `json:"table"`
}

type SmartAttribute struct {
	Id         int                 `json:"id"`
	Name       string              `json:"name"`
	Value      int                 `json:"value"`
	Worst      int                 `json:"worst"`
	Thresh     int                 `json:"thresh"`
	WhenFailed string              `json:"when_failed"`
	Flags      SmartAttributeFlags `json:"flags"`
	Raw        SmartAttributeRaw   `json:"raw"`
}

type SmartAttributeFlags struct {
	Value         int    `json:"value"`
	String        string `json:"string"`
	Prefailure    bool   `json:"prefailure"`
	UpdatedOnline bool   `json:"updated_online"`
	Performance   bool   `json:"performance"`
	ErrorRate     bool   `json:"error_rate"`
	EventCount    bool   `json:"event_count"`
	AutoKeep      bool   `json:"auto_keep"`
}

type SmartAttributeRaw struct {
	Value  int    `json:"value"`
	String string `json:"string"`
}

type AtaSctStatus struct {
	FormatVersion int         `json:"format_version"`
	SctVersion    int         `json:"sct_version"`
	DeviceState   DeviceState `json:"device_state"`
	Temperature   Temperature `json:"temperature"`
	SmartStatus   SmartStatus `json:"smart_status"`
}

type DeviceState struct {
	Value  int    `json:"value"`
	String string `json:"string"`
}

func ParseSmartctlJson(data []byte) (response Response, err error) {
	err = json.Unmarshal(data, &response)
	return
}
