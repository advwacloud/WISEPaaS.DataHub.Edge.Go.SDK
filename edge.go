package agent

import (
	"time"
)

// EdgeAgentOptions ...
type EdgeAgentOptions struct {
	ReconnectInterval int // second
	ScadaID           string
	DeviceID          string
	Type              byte
	HeartBeatInterval int
	DataRecover       bool
	ConnectType       string
	UseSecure         bool
	MQTT              *MQTTOptions
	DCCS              *DCCSOptions
}

// MQTTOptions ...
type MQTTOptions struct {
	HostName     string
	Port         int
	UserName     string
	Password     string
	ProtocalType string
}

// DCCSOptions ...
type DCCSOptions struct {
	URL string
	Key string
}

// DeviceStatus ...
type DeviceStatus struct {
	ID     string
	Status byte
}

// EdgeConfig ...
type EdgeConfig struct {
	Scada ScadaConfig
}

// EdgeData ...
type EdgeData struct {
	TagList   []EdgeTag
	Timestamp time.Time
}

// EdgeTag ...
type EdgeTag struct {
	DeviceID string
	TagName  string
	Value    interface{}
}

// EdgeDeviceStatus ...
type EdgeDeviceStatus struct {
	DeviceList []DeviceStatus
	Timestamp  time.Time
}

// ScadaConfig ...
type ScadaConfig struct {
	primaryIP   interface{}
	backupIP    interface{}
	primaryPort interface{}
	backupPort  interface{}
	scadaType   interface{}
	DeviceList  []DeviceConfig
}

// DeviceConfig ...
type DeviceConfig struct {
	id              interface{}
	name            interface{}
	comPortNumber   interface{}
	deviceType      interface{}
	description     interface{}
	ip              interface{}
	port            interface{}
	AnalogTagList   []AnalogTagConfig
	DiscreteTagList []DiscreteTagConfig
	TextTagList     []TextTagConfig
}

// AnalogTagConfig ...
type AnalogTagConfig struct {
	name                  interface{}
	description           interface{}
	readOnly              interface{}
	arraySize             interface{}
	spanHigh              interface{}
	spanLow               interface{}
	engineerUnit          interface{}
	integerDisplayFormat  interface{}
	fractionDisplayFormat interface{}
}

// DiscreteTagConfig ...
type DiscreteTagConfig struct {
	name        interface{}
	description interface{}
	readOnly    interface{}
	arraySize   interface{}
	state0      interface{}
	state1      interface{}
	state2      interface{}
	state3      interface{}
	state4      interface{}
	state5      interface{}
	state6      interface{}
	state7      interface{}
}

// TextTagConfig ...
type TextTagConfig struct {
	name        interface{}
	description interface{}
	readOnly    interface{}
	arraySize   interface{}
}

// MessageReceivedEventArgs ...
type MessageReceivedEventArgs struct {
	Type    byte
	Message interface{}
}

// ConfigAckMessage ...
type ConfigAckMessage struct {
	Result bool
}

// WriteDataMessage ...
type WriteDataMessage struct {
	DeviceList []Device
	Timestamp  time.Time
}

// Device ...
type Device struct {
	ID      string
	TagList []Tag
}

// Tag ...
type Tag struct {
	Name  string
	Value interface{}
}

// TimeSyncMessage ...
type TimeSyncMessage struct {
	UTCTime time.Time
}

// NewEdgeAgentOptions will create a new EdgeAgentOption with some new values
//	ReconnectInterval: 1
//	Type: EdgeType["Gateway"]
//	HeartBeatInterval: HeartBeatInterval
//	DataRecover: true,
//	ConnectType: ConnectType["DCCS"],
//	UseSecure: false,
//	MQTT.Port: 1883
//	MQTT.ProtocalType: Protocol["TCP"]
func NewEdgeAgentOptions() *EdgeAgentOptions {
	options := &EdgeAgentOptions{
		ReconnectInterval: 1,
		ScadaID:           "",
		DeviceID:          "",
		Type:              EdgeType["Gateway"],
		HeartBeatInterval: HeartBeatInterval,
		DataRecover:       true,
		ConnectType:       ConnectType["DCCS"],
		UseSecure:         false,
		MQTT: &MQTTOptions{
			HostName:     "",
			Port:         1883,
			UserName:     "",
			Password:     "",
			ProtocalType: Protocol["TCP"],
		},
		DCCS: &DCCSOptions{
			URL: "https://api-dccs.wise-paas.com/",
			Key: "0c053cf0329e0100c5255cfdd55defcz",
		},
	}
	return options
}

// NewScadaConfig ...
func NewScadaConfig() ScadaConfig {
	return ScadaConfig{}
}

// NewDeviceConfig ...
func NewDeviceConfig(ID string) DeviceConfig {
	return DeviceConfig{
		id: ID,
	}
}

// NewAnaglogTagConfig ...
func NewAnaglogTagConfig(name string) AnalogTagConfig {
	return AnalogTagConfig{
		name: name,
	}
}

// NewDiscreteTagConfig ...
func NewDiscreteTagConfig(name string) DiscreteTagConfig {
	return DiscreteTagConfig{
		name: name,
	}
}

// NewTextTagConfig ...
func NewTextTagConfig(name string) TextTagConfig {
	return TextTagConfig{
		name: name,
	}
}

func getWriteDataMessageFromCmdMessage(data interface{}) WriteDataMessage {
	m := data.(map[string]interface{})
	message := WriteDataMessage{
		Timestamp: time.Now(),
	}
	for device, value := range m {
		d := Device{
			ID: device,
		}
		tagList := value.(map[string]interface{})
		for tag, v := range tagList {
			t := Tag{
				Name:  tag,
				Value: v,
			}
			d.TagList = append(d.TagList, t)
		}
		message.DeviceList = append(message.DeviceList, d)
	}

	return message
}

func getTimeSyncMessageFromCmdMessage(utc int) TimeSyncMessage {
	ts := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC)
	duration := time.Duration(utc) * time.Second
	message := TimeSyncMessage{
		UTCTime: ts.Add(duration),
	}
	return message
}

func (o *MQTTOptions) isValid() bool {
	return !(o.HostName == "" || o.Port == 0 || o.ProtocalType == "")
}

func (o *DCCSOptions) isValid() bool {
	return !(o.URL == "" || o.Key == "")
}

// SetPrimaryIP ...
func (config *ScadaConfig) SetPrimaryIP(ip string) {
	config.primaryIP = ip
}

// SetBackupIP ...
func (config *ScadaConfig) SetBackupIP(ip string) {
	config.backupIP = ip
}

// SetPrimaryPort ...
func (config *ScadaConfig) SetPrimaryPort(port int) {
	config.primaryPort = port
}

// SetBackupPort ...
func (config *ScadaConfig) SetBackupPort(port int) {
	config.backupPort = port
}

// SetType ...
func (config *ScadaConfig) SetType(t byte) {
	config.scadaType = t
}

// SetName ...
func (config *DeviceConfig) SetName(name string) {
	config.name = name
}

// SetComPortNumber ...
func (config *DeviceConfig) SetComPortNumber(num int) {
	config.comPortNumber = num
}

// SetType ...
func (config *DeviceConfig) SetType(deviceType string) {
	config.deviceType = deviceType
}

// SetDescription ...
func (config *DeviceConfig) SetDescription(desc string) {
	config.description = desc
}

// SetIP ...
func (config *DeviceConfig) SetIP(ip string) {
	config.ip = ip
}

// SetPort ...
func (config *DeviceConfig) SetPort(port int) {
	config.port = port
}

// SetDescription ...
func (config *AnalogTagConfig) SetDescription(desc string) {
	config.description = desc
}

// SetReadOnly ...
func (config *AnalogTagConfig) SetReadOnly(readOnly bool) {
	config.readOnly = readOnly
}

// SetArraySize ...
func (config *AnalogTagConfig) SetArraySize(num uint) {
	config.arraySize = num
}

// SetSpanHigh ...
func (config *AnalogTagConfig) SetSpanHigh(high float32) {
	config.spanHigh = high
}

// SetSpanLow ...
func (config *AnalogTagConfig) SetSpanLow(low float32) {
	config.spanLow = low
}

// SetEngineerUnit ...
func (config *AnalogTagConfig) SetEngineerUnit(unit string) {
	config.engineerUnit = unit
}

// SetIntegerDisplayFormat ...
func (config *AnalogTagConfig) SetIntegerDisplayFormat(num uint) {
	config.integerDisplayFormat = num
}

// SetFractionDisplayFormat ...
func (config *AnalogTagConfig) SetFractionDisplayFormat(num uint) {
	config.fractionDisplayFormat = num
}

// SetDescription ...
func (config *DiscreteTagConfig) SetDescription(desc string) {
	config.description = desc
}

// SetReadOnly ...
func (config *DiscreteTagConfig) SetReadOnly(readOnly bool) {
	config.readOnly = readOnly
}

// SetArraySize ...
func (config *DiscreteTagConfig) SetArraySize(num uint) {
	config.arraySize = num
}

// SetState0 ...
func (config *DiscreteTagConfig) SetState0(state string) {
	config.state0 = state
}

// SetState1 ...
func (config *DiscreteTagConfig) SetState1(state string) {
	config.state1 = state
}

// SetState2 ...
func (config *DiscreteTagConfig) SetState2(state string) {
	config.state2 = state
}

// SetState3 ...
func (config *DiscreteTagConfig) SetState3(state string) {
	config.state3 = state
}

// SetState4 ...
func (config *DiscreteTagConfig) SetState4(state string) {
	config.state4 = state
}

// SetState5 ...
func (config *DiscreteTagConfig) SetState5(state string) {
	config.state5 = state
}

// SetState6 ...
func (config *DiscreteTagConfig) SetState6(state string) {
	config.state6 = state
}

// SetState7 ...
func (config *DiscreteTagConfig) SetState7(state string) {
	config.state7 = state
}

// SetDescription ...
func (config *TextTagConfig) SetDescription(desc string) {
	config.description = desc
}

// SetReadOnly ...
func (config *TextTagConfig) SetReadOnly(readOnly bool) {
	config.readOnly = readOnly
}

// SetArraySize ...
func (config *TextTagConfig) SetArraySize(num uint) {
	config.arraySize = num
}
