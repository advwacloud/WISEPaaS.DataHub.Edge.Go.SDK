package agent

const (
	// tag key for map
	tagKeyFormat string = "%s|%s|%s"
	// HeartBeatInterval ...
	HeartBeatInterval int = 60 // second
	// defaultReadRecordCount ...
	defaultReadRecordCount int = 10
	// dataRecoverInterval ...
	dataRecoverInterval int = 3 //second
	// dataRecoverFilePath ...
	dataRecoverFilePath string = "_recover.sqlite"
	// tags conifg file path
	tagsCfgFilePath string = "_config.json"
	// limit data size
	dataMaxTagCount int = 100
)

// Action ...
var Action = map[string]byte{
	"Create":  1,
	"Update":  2,
	"Delete":  3,
	"Delsert": 4,
}

// ConnectType ...
var ConnectType = map[string]string{
	"MQTT": "MQTT",
	"DCCS": "DCCS",
}

// EdgeType ...
var EdgeType = map[string]byte{
	"Gateway": 0,
	"Device":  1,
}

// MessageType ...
var MessageType = map[string]byte{
	"WriteValue":  0,
	"WrtieConfig": 1,
	"TimeSync":    2,
	"ConfigAck":   3,
}

// mqttQoS ...
var mqttQoS = map[string]byte{
	"AtMostOnce":  0,
	"AtLeastOnce": 1,
	"ExactlyOnce": 2,
}

// Protocol ...
var Protocol = map[string]string{
	"TCP":       "tcp",
	"WebSocket": "websockets",
	"TLS":       "tls",
}

var protocolScheme = map[string]string{
	"tcp":        "tcp",
	"websockets": "ws",
	"tls":        "tls",
}

// Status ...
var Status = map[string]byte{
	"Offline": 0,
	"Online":  1,
}

// TagType ...
var TagType = map[string]byte{
	"Analog":   1,
	"Discrete": 2,
	"Text":     3,
}
