package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type tagsCfgHelper interface {
	addCfgByUploadConfig(a *agent, config *EdgeConfig) bool
	addCfgFromFile(a *agent, filePath string) bool
	writeCfgToFile(a *agent, filePath string) bool
}

type tagsCfgStruct struct{}

func newTagsCfgHelper() tagsCfgHelper {
	return &tagsCfgStruct{}
}

func (helper *tagsCfgStruct) addCfgByUploadConfig(a *agent, config *EdgeConfig) bool {
	nodeID := a.options.NodeID

	a.tagsCfgMap = make(map[string]map[string]interface{})

	for _, device := range config.Node.DeviceList {
		for _, tag := range device.AnalogTagList {
			tagKey := fmt.Sprintf(tagKeyFormat, nodeID, device.id, tag.name)

			cfg, ok := a.tagsCfgMap[tagKey]

			if !ok {
				cfg = make(map[string]interface{})
				a.tagsCfgMap[tagKey] = cfg
			}

			a.tagsCfgMap[tagKey]["fractionDisplayFormat"] = tag.fractionDisplayFormat

		}
	}

	return true
}

func (helper *tagsCfgStruct) addCfgFromFile(a *agent, filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	if err := json.Unmarshal([]byte(content), &a.tagsCfgMap); err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	return true
}

func (helper *tagsCfgStruct) writeCfgToFile(a *agent, filePath string) bool {
	jsonStr, err := json.Marshal(a.tagsCfgMap)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	err = ioutil.WriteFile("tagsCfgMap.json", []byte(jsonStr), 0644)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	return true
}
