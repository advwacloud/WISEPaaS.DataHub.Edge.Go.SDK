package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	UUID "github.com/google/uuid"
)

// Agent ...
type Agent interface {
	IsConnected() bool
	Connect() error
	Disconnect()
	SetOnConnectHandler(onConn OnConnectHandler)
	SetOnDisconnectHandler(onDisconn OnDisconnectHandler)
	SetOnMessageReceiveHandler(onMessageReceive OnMessageReceiveHandler)
	UploadConfig(action byte, edgeConfig EdgeConfig) bool
	SendDeviceStatus(status EdgeDeviceStatus) bool
	SendData(data EdgeData) bool
}

// Agent ...
type agent struct {
	options           EdgeAgentOptions
	client            MQTT.Client // interface
	heartbeatTimer    chan bool
	dataRecoverTimer  chan bool
	dataRecoverHelper DataRecoverHelper
	cfgCache          configMessage
	OnConnect         OnConnectHandler
	OnDisconnect      OnDisconnectHandler
	OnMessageReceive  OnMessageReceiveHandler
}

// OnConnectHandler ...
type OnConnectHandler func(Agent)

// OnDisconnectHandler ...
type OnDisconnectHandler func(Agent)

// OnMessageReceiveHandler ...
type OnMessageReceiveHandler func(MessageReceivedEventArgs)

// NewAgent ...
func NewAgent(options *EdgeAgentOptions) Agent {
	a := &agent{
		options:           *options,
		client:            nil,
		heartbeatTimer:    nil,
		dataRecoverTimer:  nil,
		dataRecoverHelper: nil,
		cfgCache:          configMessage{},
		OnConnect:         func(a Agent) {},
		OnDisconnect:      func(a Agent) {},
		OnMessageReceive:  func(res MessageReceivedEventArgs) {},
	}
	if options.DataRecover {
		a.dataRecoverHelper = NewDataRecoverHelper(dataRecoverFilePath)
	}

	// add cfg to memory from disk
	helper := newTagsCfgHelper()
	helper.getCfgFromFile(a, tagsCfgFilePath)

	return a
}

// IsConnected ...
func (a *agent) IsConnected() bool {
	if a.client == nil {
		return false
	}
	return a.client.IsConnectionOpen()
}

// Connect ...
func (a *agent) Connect() error {

	if a.IsConnected() {
		return nil
	}
	if a.options.ConnectType == ConnectType["DCCS"] {
		if !a.options.DCCS.isValid() {
			return errors.New("DCCS options is invalid")
		}
		error := a.getCredentailFromDCCS()
		if error != nil {
			fmt.Println(error)
			return error
		}
	}
	if !a.options.MQTT.isValid() {
		return errors.New("MQTT options is invalid")
	}

	clientOptions, _ := a.newClientOptions()
	a.client = MQTT.NewClient(clientOptions)
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Disconnect ...
func (a *agent) Disconnect() {
	if !a.client.IsConnected() {
		return
	}

	/* Send Disconnect message */
	topic := fmt.Sprintf(mqttTopic["DeviceConnTopic"], a.options.NodeID, a.options.DeviceID)
	if a.options.Type == EdgeType["GateWay"] {
		topic = fmt.Sprintf(mqttTopic["NodeConnTopic"], a.options.NodeID)
	}
	payload := newDisconnectMessage().getPayload()
	if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], true, payload); token.Wait() && token.Error() != nil {
		fmt.Println("token error in Disconnect: ", token.Error())
	}

	go a.handleDisconnect()
	a.client.Disconnect(0)
}

func (a *agent) UploadConfig(action byte, config EdgeConfig) bool {
	if !a.IsConnected() {
		return false
	}
	nodeID := a.options.NodeID

	var payload configMessage
	var result = false
	switch action {
	case Action["Create"]:
		result, payload = convertCreateorUpdateConfig(action, nodeID, config, a.options.HeartBeatInterval)
	case Action["Update"]:
		result, payload = convertCreateorUpdateConfig(action, nodeID, config, a.options.HeartBeatInterval)
	case Action["Delete"]:
		result, payload = convertDeleteConfig(action, nodeID, config)
	case Action["Delsert"]:
		result, payload = convertCreateorUpdateConfig(action, nodeID, config, a.options.HeartBeatInterval)
	default:
		result = false
	}

	if action != Action["Delete"] {
		helper := newTagsCfgHelper()

		// add config to memory
		helper.addCfgToMemory(a, payload)

		// write config to disk
		helper.addCfgToFile(a, tagsCfgFilePath)
	}

	if result {
		topic := fmt.Sprintf(mqttTopic["ConfigTopic"], a.options.NodeID)
		if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], true, payload.getPayload()); token.Wait() && token.Error() != nil {
			fmt.Println("token error in UploadConfig: ", token.Error())
			result = false
		}
	}
	return result
}

func (a *agent) SendDeviceStatus(statuses EdgeDeviceStatus) bool {
	if !(a.IsConnected()) {
		return false
	}
	msg := newStatusMessage()
	msg.Ts = statuses.Timestamp.Format(time.RFC3339)
	for _, status := range statuses.DeviceList {
		msg.D.Dev[status.ID] = status.Status
	}
	payload := msg.getPayload()
	topic := fmt.Sprintf(mqttTopic["NodeConnTopic"], a.options.NodeID)
	if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], true, payload); token.Wait() && token.Error() != nil {
		fmt.Println("token error in SendDviceStatus: ", token.Error())
		return false
	}
	return true
}

func (a *agent) SendData(data EdgeData) bool {
	result, payloads := convertTagValue(data, a)
	topic := fmt.Sprintf(mqttTopic["DataTopic"], a.options.NodeID)
	if !a.IsConnected() {
		for _, payload := range payloads {
			if a.dataRecoverHelper != nil {
				a.dataRecoverHelper.Write(payload)
			}
		}
		result = false
	} else {
		for _, payload := range payloads {
			if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], true, payload); token.Wait() && token.Error() != nil {
				fmt.Println("token error in SendData: ", token.Error())
				if a.dataRecoverHelper != nil {
					a.dataRecoverHelper.Write(payload)
				}
				result = false
			}
		}
	}
	return result
}

func (a *agent) getCredentailFromDCCS() error {
	url := a.options.DCCS.URL
	if url[len(url)-1:] == "/" {
		a.options.DCCS.URL = url[:len(url)-1]
	}
	url = fmt.Sprintf("%s/v1/serviceCredentials/%s", a.options.DCCS.URL, a.options.DCCS.Key)
	res, error := http.Get(url)
	if error != nil {
		return error
	}

	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return error
	}

	var response struct {
		ServiceName string
		ServiceHost string
		Credential  struct {
			Password  string
			Username  string
			Protocols map[string]struct {
				Ssl      bool
				Username string
				Password string
				Port     int
			}
		}
	}
	error = json.Unmarshal([]byte(body), &response)
	if error != nil {
		return error
	}

	a.options.MQTT.HostName = response.ServiceHost
	if a.options.UseSecure {
		a.options.MQTT.Port = response.Credential.Protocols["mqtt+ssl"].Port
		a.options.MQTT.UserName = response.Credential.Protocols["mqtt+ssl"].Username
		a.options.MQTT.Password = response.Credential.Protocols["mqtt+ssl"].Password
		a.options.MQTT.ProtocalType = Protocol["TLS"]
	} else {
		a.options.MQTT.Port = response.Credential.Protocols["mqtt"].Port
		a.options.MQTT.UserName = response.Credential.Protocols["mqtt"].Username
		a.options.MQTT.Password = response.Credential.Protocols["mqtt"].Password
	}
	return nil
}

func (a *agent) newClientOptions() (*MQTT.ClientOptions, error) {
	clientOptions := MQTT.NewClientOptions()
	schema := protocolScheme[Protocol["TCP"]]
	// Enable Debug
	// MQTT.DEBUG = log.New(os.Stdout, "[Debug] ", 0)
	// MQTT.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	// MQTT.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	// MQTT.WARN = log.New(os.Stdout, "[WARN]  ", 0)

	if a.options.MQTT.ProtocalType == Protocol["WebSocket"] {
		schema = protocolScheme[Protocol["WebSocket"]]
	}
	if a.options.MQTT.ProtocalType == Protocol["TLS"] {
		schema = protocolScheme[Protocol["TLS"]]
	}

	server := fmt.Sprintf("%s://%s:%d", schema, a.options.MQTT.HostName, a.options.MQTT.Port)
	clientOptions.AddBroker(server)
	uuid := UUID.New()
	clientOptions.SetClientID(fmt.Sprintf("EdgeAgent_%s", uuid))
	clientOptions.SetAutoReconnect(true)
	clientOptions.SetConnectRetry(true)
	clientOptions.SetConnectRetryInterval(time.Duration(a.options.ReconnectInterval) * time.Second)
	clientOptions.SetCleanSession(false)
	clientOptions.SetPassword(a.options.MQTT.Password)
	clientOptions.SetUsername(a.options.MQTT.UserName)
	clientOptions.SetMaxReconnectInterval(time.Duration(a.options.ReconnectInterval) * time.Second)
	topic := fmt.Sprintf(mqttTopic["NodeConnTopic"], a.options.NodeID)
	payload := newWillMessage().getPayload()
	clientOptions.SetWill(topic, payload, mqttQoS["AtLeastOnce"], true)

	clientOptions.SetOnConnectHandler(a.handleOnConnect)
	clientOptions.SetConnectionLostHandler(func(c MQTT.Client, err error) {
		fmt.Println("Reconnecting...")
		if err != nil {
			fmt.Println(err)
		}
		if a.options.ConnectType == ConnectType["DCCS"] {
			error := a.getCredentailFromDCCS()
			if error != nil {
				fmt.Println(err)
			}
		}
		go a.OnDisconnect(a)
	})
	return clientOptions, nil
}

func (a *agent) SetOnConnectHandler(onConn OnConnectHandler) {
	a.OnConnect = onConn
}

func (a *agent) SetOnDisconnectHandler(onDisconn OnDisconnectHandler) {
	a.OnDisconnect = onDisconn
}

func (a *agent) SetOnMessageReceiveHandler(onMessageReceive OnMessageReceiveHandler) {
	a.OnMessageReceive = onMessageReceive
}

func (a *agent) handleOnConnect(c MQTT.Client) {
	/* subscribe */
	cmdTopic := fmt.Sprintf(mqttTopic["DeviceCmdTopic"], a.options.NodeID, a.options.DeviceID)
	if a.options.Type == EdgeType["Gateway"] {
		cmdTopic = fmt.Sprintf(mqttTopic["NodeCmdTopic"], a.options.NodeID)
	}
	if token := a.client.Subscribe(cmdTopic, mqttQoS["AtLeastOnce"], a.handleCmdReceive); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	ackTopic := fmt.Sprintf(mqttTopic["AckTopic"], a.options.NodeID)
	if token := a.client.Subscribe(ackTopic, mqttQoS["AtLeastOnce"], a.handleAckReceive); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	/* Send connect Message */
	topic := fmt.Sprintf(mqttTopic["DeviceConnTopic"], a.options.NodeID, a.options.DeviceID)
	if a.options.Type == EdgeType["GateWay"] {
		topic = fmt.Sprintf(mqttTopic["NodeConnTopic"], a.options.NodeID)
	}
	payload := newConnMessage().getPayload()
	if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], true, payload); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	/* heartbeat */
	if a.options.HeartBeatInterval > 0 && a.heartbeatTimer == nil {
		interval := a.options.HeartBeatInterval
		a.heartbeatTimer = setInterval(a.sendHeartBeat, interval, true)
	}

	/* Recover */
	if a.options.DataRecover && a.dataRecoverTimer == nil {
		a.dataRecoverTimer = setInterval(a.sendRecover, dataRecoverInterval, true)
	}

	go a.OnConnect(a)
}

func (a *agent) handleDisconnect() {
	for a.client.IsConnectionOpen() {
	}
	fmt.Println("Disconnected...")
	a.client = nil
	a.heartbeatTimer <- false
	a.heartbeatTimer = nil
	a.dataRecoverTimer <- false
	a.dataRecoverTimer = nil
	go a.OnDisconnect(a)
}

func (a *agent) handleCmdReceive(c MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	if !isJSON(payload) {
		fmt.Println("Invalid JSON format")
		return
	}
	var data cmdMessage
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		fmt.Println("Cmd decode failed:", err)
	}
	var message interface{}
	var argType byte
	switch data.D.Cmd {
	case "WV":
		argType = MessageType["WriteValue"]
		message = getWriteDataMessageFromCmdMessage(data.D.Val, data.Ts)
	case "TSyn":
		argType = MessageType["TimeSync"]
		message = getTimeSyncMessageFromCmdMessage(data.D.UTC)
	default:
		// fmt.Println("Message format is invalid")
		return
	}
	res := MessageReceivedEventArgs{
		Type:    argType,
		Message: message,
	}
	go a.OnMessageReceive(res)
}

func (a *agent) handleAckReceive(c MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	if !isJSON(payload) {
		fmt.Println("Invalid JSON format")
		return
	}
	var data ackConfigMessage
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		fmt.Println("Ack decode failed:", err)
		return
	}
	val, ok := data.D.Cfg.(float64)
	if data.D.Cfg == nil || !ok {
		// fmt.Println("Message format is invalid")
		return
	}
	var result = false
	if val > 0 {
		result = true
	}
	message := ConfigAckMessage{
		Result: result,
	}
	res := MessageReceivedEventArgs{
		Type:    MessageType["ConfigAck"],
		Message: message,
	}
	go a.OnMessageReceive(res)
}

func (a *agent) sendHeartBeat() {
	if !a.IsConnected() {
		return
	}
	topic := fmt.Sprintf(mqttTopic["DeviceConnTopic"], a.options.NodeID, a.options.DeviceID)
	if a.options.Type == EdgeType["GateWay"] {
		topic = fmt.Sprintf(mqttTopic["NodeConnTopic"], a.options.NodeID)
	}
	payload := newHeartBeatMessage().getPayload()
	if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], true, payload); token.Wait() && token.Error() != nil {
		fmt.Println("token error in sendHeartBeat: ", token.Error())
	}
}

func (a *agent) sendRecover() {
	if !a.IsConnected() {
		return
	}
	if a.dataRecoverHelper == nil {
		return
	}
	helper := a.dataRecoverHelper

	if !helper.IsDataExist() {
		return
	}
	messages := helper.Read(defaultReadRecordCount)
	topic := fmt.Sprintf(mqttTopic["DataTopic"], a.options.NodeID)
	for _, message := range messages {
		if !a.IsConnected() {
			helper.Write(message)
			continue
		}
		if token := a.client.Publish(topic, mqttQoS["AtLeastOnce"], false, message); token.Wait() && token.Error() != nil {
			fmt.Println("token error in sendRecover: ", token.Error())
			helper.Write(message)
		}
	}
}
