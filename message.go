package agent

import (
	"encoding/json"
	"time"
)

// Message ...
type Message interface {
	getPayload() string
}

type message struct {
	Ts string      `json:"ts"`
	D  interface{} `json:"d"`
}

type willData struct {
	UeD byte
}

type heartbeatData struct {
	Hbt byte
}

type connectData struct {
	Con byte
}

type disconnectData struct {
	DsC byte
}

type configMessage struct {
	Ts string     `json:"ts"`
	D  configData `json:"d"`
}

type configData struct {
	Action byte
	Scada  map[string]interface{}
}

type statusMessage struct {
	Ts string     `json:"ts"`
	D  statusData `json:"d"`
}

type statusData struct {
	Dev map[string]byte
}

type tagValue struct {
	Ts string                 `json:"ts"`
	D  map[string]interface{} `json:"d"`
}

// newWillMessage ...
func newWillMessage() Message {
	msg := &message{
		D: willData{
			UeD: 1,
		},
		Ts: time.Now().UTC().Format(time.RFC3339),
	}
	return msg
}

// newHeartBeatMessage ...
func newHeartBeatMessage() Message {
	message := &message{
		D: heartbeatData{
			Hbt: 1,
		},
		Ts: time.Now().UTC().Format(time.RFC3339),
	}
	return message
}

func newConnMessage() Message {
	message := &message{
		D: connectData{
			Con: 1,
		}, Ts: time.Now().UTC().Format(time.RFC3339),
	}
	return message
}

func newDisconnectMessage() Message {
	message := &message{
		D: disconnectData{
			DsC: 1,
		}, Ts: time.Now().UTC().Format(time.RFC3339),
	}
	return message
}

// newConfigData ...
func newConfigData(action byte) configMessage {
	return configMessage{
		D: configData{
			Action: action,
			Scada:  make(map[string]interface{}),
		},
		Ts: time.Now().UTC().Format(time.RFC3339),
	}
}

func newStatusMessage() statusMessage {
	return statusMessage{
		D: statusData{
			Dev: make(map[string]byte),
		},
		Ts: time.Now().UTC().Format(time.RFC3339),
	}
}

func newTagValue(ts time.Time) tagValue {
	currentTimeData := time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), time.Local)

	t := tagValue{
		D:  make(map[string]interface{}),
		Ts: currentTimeData.UTC().Format(time.RFC3339Nano),
	}
	return t
}

func (m *message) getPayload() string {
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (m *configMessage) getPayload() string {
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (m *statusMessage) getPayload() string {
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (m *tagValue) getPayload() string {
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(j)
}

type ackConfigMessage struct {
	Ts string
	D  struct {
		Cfg interface{}
	}
}

type cmdMessage struct {
	Ts string
	D  struct {
		Cmd string
		Val interface{}
		UTC int
	}
}
