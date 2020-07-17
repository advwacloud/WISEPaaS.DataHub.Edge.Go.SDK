package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type tagsCfgHelper interface {
	addCfgToMemory(a *agent, config configMessage) bool
	getCfgFromFile(a *agent, filePath string) bool
	addCfgToFile(a *agent, filePath string) bool
}

type tagsCfgStruct struct{}

func newTagsCfgHelper() tagsCfgHelper {
	return &tagsCfgStruct{}
}

func (helper *tagsCfgStruct) addCfgToMemory(a *agent, config configMessage) bool {
	a.cfgCache = config
	return true
}

func (helper *tagsCfgStruct) getCfgFromFile(a *agent, filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	if err := json.Unmarshal([]byte(content), &a.cfgCache); err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	return true
}

func (helper *tagsCfgStruct) addCfgToFile(a *agent, filePath string) bool {
	jsonStr, err := json.Marshal(a.cfgCache)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	err = ioutil.WriteFile(filePath, []byte(jsonStr), 0644)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	return true
}
