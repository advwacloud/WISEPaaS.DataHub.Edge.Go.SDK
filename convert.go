package agent

import (
	"fmt"
	"sort"
	"strconv"
)

func convertCreateorUpdateConfig(action byte, nodeID string, config EdgeConfig, heartbeat int) (bool, string) {
	message := newConfigData(action)
	message.D.Scada[nodeID] = make(map[string]interface{})
	message.D.Scada[nodeID] = convertNodeConfig(config.Node, heartbeat)
	return true, message.getPayload()
}

func convertDeleteConfig(action byte, nodeID string, config EdgeConfig) (bool, string) {
	message := newConfigData(action)
	node := config.Node
	s := make(map[string]interface{})
	for _, device := range node.DeviceList {
		if s["Device"] == nil {
			s["Device"] = make(map[string]interface{})
		}
		id, d := convertDeviceConfig(device)
		s["Device"].(map[string]interface{})[id] = d
	}
	message.D.Scada[nodeID] = make(map[string]interface{})
	message.D.Scada[nodeID] = s
	return true, message.getPayload()
}

func convertNodeConfig(config NodeConfig, heartbeat int) map[string]interface{} {
	s := make(map[string]interface{})
	s["Hbt"] = heartbeat
	if config.primaryIP != nil {
		s["PIP"] = config.primaryIP
	}
	if config.backupIP != nil {
		s["BIP"] = config.backupIP
	}
	if config.primaryPort != nil {
		s["PPort"] = config.primaryPort
	}
	if config.backupPort != nil {
		s["BPort"] = config.backupPort
	}
	if config.nodeType != nil {
		s["Type"] = config.nodeType
	}

	for _, device := range config.DeviceList {
		if s["Device"] == nil {
			s["Device"] = make(map[string]interface{})
		}
		id, d := convertDeviceConfig(device)
		s["Device"].(map[string]interface{})[id] = d
	}

	return s
}

func convertDeviceConfig(device DeviceConfig) (string, map[string]interface{}) {
	d := make(map[string]interface{})
	if device.name != nil {
		d["Name"] = device.name
	}
	if device.comPortNumber != nil {
		d["PNbr"] = device.comPortNumber
	}
	if device.deviceType != nil {
		d["Type"] = device.deviceType
	}
	if device.description != nil {
		d["Desc"] = device.description
	}
	if device.ip != nil {
		d["IP"] = device.ip
	}
	if device.port != nil {
		d["Port"] = device.port
	}
	if device.retentionPolicyName != nil {
		d["RP"] = device.retentionPolicyName
	}

	for _, tag := range device.AnalogTagList {
		if d["Tag"] == nil {
			d["Tag"] = make(map[string]interface{})
		}
		name, t := convertAnalogTagConfig(tag)
		d["Tag"].(map[string]interface{})[name] = t
	}
	for _, tag := range device.DiscreteTagList {
		if d["Tag"] == nil {
			d["Tag"] = make(map[string]interface{})
		}
		name, t := convertDiscreteTagConfig(tag)
		d["Tag"].(map[string]interface{})[name] = t
	}
	for _, tag := range device.TextTagList {
		if d["Tag"] == nil {
			d["Tag"] = make(map[string]interface{})
		}
		name, t := convertTextTagConfig(tag)
		d["Tag"].(map[string]interface{})[name] = t
	}
	return device.id.(string), d
}

func convertAnalogTagConfig(tag AnalogTagConfig) (string, map[string]interface{}) {
	t := make(map[string]interface{})
	t["Type"] = TagType["Analog"]
	if tag.description != nil {
		t["Desc"] = tag.description
	}
	if tag.readOnly != nil {
		t["RO"] = 0
		if tag.readOnly.(bool) {
			t["RO"] = 1
		}
	}
	if tag.arraySize != nil {
		t["Ary"] = tag.arraySize
	}
	if tag.spanHigh != nil {
		t["SH"] = tag.spanHigh
	}
	if tag.spanLow != nil {
		t["SL"] = tag.spanLow
	}
	if tag.engineerUnit != nil {
		t["EU"] = tag.engineerUnit
	}
	if tag.integerDisplayFormat != nil {
		t["IDF"] = tag.integerDisplayFormat
	}
	if tag.fractionDisplayFormat != nil {
		t["FDF"] = tag.fractionDisplayFormat
	}
	return tag.name.(string), t
}

func convertDiscreteTagConfig(tag DiscreteTagConfig) (string, map[string]interface{}) {
	t := make(map[string]interface{})
	t["Type"] = TagType["Discrete"]
	if tag.description != nil {
		t["Desc"] = tag.description
	}
	if tag.readOnly != nil {
		t["RO"] = 0
		if tag.readOnly.(bool) {
			t["RO"] = 1
		}
	}
	if tag.arraySize != nil {
		t["Ary"] = tag.arraySize
	}
	if tag.state0 != nil {
		t["S0"] = tag.state0
	}
	if tag.state1 != nil {
		t["S1"] = tag.state1
	}
	if tag.state2 != nil {
		t["S2"] = tag.state2
	}
	if tag.state3 != nil {
		t["S3"] = tag.state3
	}
	if tag.state4 != nil {
		t["S4"] = tag.state4
	}
	if tag.state5 != nil {
		t["S5"] = tag.state5
	}
	if tag.state6 != nil {
		t["S6"] = tag.state6
	}
	if tag.state7 != nil {
		t["S7"] = tag.state7
	}
	return tag.name.(string), t
}

func convertTextTagConfig(tag TextTagConfig) (string, map[string]interface{}) {
	t := make(map[string]interface{})
	t["Type"] = TagType["Text"]
	if tag.description != nil {
		t["Desc"] = tag.description
	}
	if tag.readOnly != nil {
		t["RO"] = 0
		if tag.readOnly.(bool) {
			t["RO"] = 1
		}
	}
	if tag.arraySize != nil {
		t["Ary"] = tag.arraySize
	}
	return tag.name.(string), t
}

func convertTagValue(data EdgeData, a *agent) (bool, []string) {
	count := 0
	list := data.TagList
	var messages []string
	msg := newTagValue()

	sort.Slice(list[:], func(i, j int) bool {
		return list[i].DeviceID < list[j].DeviceID
	})

	for _, tag := range list {

		if msg.D[tag.DeviceID] == nil {
			msg.D[tag.DeviceID] = make(map[string]interface{})
		}

		tagKey := fmt.Sprintf(tagKeyFormat, a.options.NodeID, tag.DeviceID, tag.TagName)
		fractionDisplayFormat, ok := a.tagsCfgMap[tagKey]["fractionDisplayFormat"]

		if ok == true {
			// Round down tag value to the specified digit
			convertVal := roundDownByFDF(tag.Value, fractionDisplayFormat)
			msg.D[tag.DeviceID].(map[string]interface{})[tag.TagName] = convertVal
		} else {
			msg.D[tag.DeviceID].(map[string]interface{})[tag.TagName] = tag.Value
		}

		count++
		if count == dataMaxTagCount {
			messages = append(messages, msg.getPayload())
			msg = newTagValue()
		}
	}
	messages = append(messages, msg.getPayload())
	return true, messages
}

func roundDownByFDF(originVal interface{}, fractionDisplayFormat interface{}) float64 {
	valFormat := "%." + fmt.Sprint(fractionDisplayFormat) + "f"
	valStr := fmt.Sprintf(valFormat, originVal)

	finalVal, err := strconv.ParseFloat(valStr, 64)

	if err != nil {
		return originVal.(float64)
	}

	return finalVal
}
