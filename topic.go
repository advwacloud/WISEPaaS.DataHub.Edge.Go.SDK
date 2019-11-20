package agent

// Topic ...
var mqttTopic = map[string]string{
	"ConfigTopic":     "/wisepaas/scada/%s/cfg",
	"DataTopic":       "/wisepaas/scada/%s/data",
	"ScadaConnTopic":  "/wisepaas/scada/%s/conn",
	"DeviceConnTopic": "/wisepaas/scada/%s/%s/conn",
	"ScadaCmdTopic":   "/wisepaas/scada/%s/cmd",
	"DeviceCmdTopic":  "/wisepaas/scada/%s/%s/cmd",
	"AckTopic":        "/wisepaas/scada/%s/ack",
	"CfgAckTopic":     "/wisepaas/scada/%s/cfgack",
}
