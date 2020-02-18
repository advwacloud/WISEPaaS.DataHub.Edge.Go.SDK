package agent

// Topic ...
var mqttTopic = map[string]string{
	"ConfigTopic":     "/wisepaas/scada/%s/cfg",
	"DataTopic":       "/wisepaas/scada/%s/data",
	"NodeConnTopic":  "/wisepaas/scada/%s/conn",
	"DeviceConnTopic": "/wisepaas/scada/%s/%s/conn",
	"NodeCmdTopic":   "/wisepaas/scada/%s/cmd",
	"DeviceCmdTopic":  "/wisepaas/scada/%s/%s/cmd",
	"AckTopic":        "/wisepaas/scada/%s/ack",
	"CfgAckTopic":     "/wisepaas/scada/%s/cfgack",
}
