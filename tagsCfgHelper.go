package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
	_ "github.com/mattn/go-sqlite3"
)

type tagsCfgHelper interface {
	getCfgFromFile(a *agent, filePath string) bool
	addCfgToFile(a *agent, filePath string) bool
	overwriteCfgCache(a *agent, config EdgeConfig)
	updateCfgCache(a *agent, config EdgeConfig)
	deleteCfgCache(a *agent, config EdgeConfig)
}

type tagsCfgStruct struct{}

type configCache struct {
	// TODO: add new entry allowFormat for FDF formating logic
	deviceMap      map[string]map[string]map[string]interface{}
	recentValueMap map[string]interface{}
}

func newConfigCache() configCache {
	return configCache{
		deviceMap:      make(map[string]map[string]map[string]interface{}),
		recentValueMap: make(map[string]interface{}),
	}
}

func newTagsCfgHelper() tagsCfgHelper {
	return &tagsCfgStruct{}
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

	if err := json.Unmarshal([]byte(content), &a.cfgCache.deviceMap); err != nil {
		fmt.Printf("%s", err.Error())
		return false
	}

	return true
}

func (helper *tagsCfgStruct) addCfgToFile(a *agent, filePath string) bool {
	jsonStr, err := json.MarshalIndent(a.cfgCache.deviceMap, "", "\t")

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

func (helper *tagsCfgStruct) overwriteCfgCache(a *agent, config EdgeConfig) {
	newCfgCache := convertEdgeConfigToConfigCache(a, config, false)
	a.cfgCache = newCfgCache
}

func (helper *tagsCfgStruct) updateCfgCache(a *agent, config EdgeConfig) {
	newCfgCache := convertEdgeConfigToConfigCache(a, config, true)
	if err := mergo.Merge(&a.cfgCache.deviceMap, newCfgCache.deviceMap, mergo.WithOverride); err != nil {
		fmt.Println(err)
	}
}

func (helper *tagsCfgStruct) deleteCfgCache(a *agent, config EdgeConfig) {

	for _, deviceConfig := range config.Node.DeviceList {

		deviceID := deviceConfig.id.(string)

		if len(deviceConfig.AnalogTagList)+len(deviceConfig.DiscreteTagList)+len(deviceConfig.TextTagList) == 0 {

			delete(a.cfgCache.deviceMap, deviceID) // contain no tagData -> delete whole device

		} else {

			// TODO: need to log if delete target is not exist? it'll not throw any error tho
			for _, analogTag := range deviceConfig.AnalogTagList {
				delete(a.cfgCache.deviceMap[deviceID], analogTag.name.(string))
			}
			for _, discreteTag := range deviceConfig.DiscreteTagList {
				delete(a.cfgCache.deviceMap[deviceID], discreteTag.name.(string))
			}
			for _, textTag := range deviceConfig.TextTagList {
				delete(a.cfgCache.deviceMap[deviceID], textTag.name.(string))
			}
		}
	}
}

func convertEdgeConfigToConfigCache(a *agent, config EdgeConfig, isUpdate bool) configCache {
	newCfgCache := newConfigCache()
	newDeviceMap := newCfgCache.deviceMap

	for _, deviceConfig := range config.Node.DeviceList {

		deviceID := deviceConfig.id.(string)
		if _, exist := a.cfgCache.deviceMap[deviceID]; !exist && isUpdate {
			continue // skip unexist device when the action is 'Update'
		}

		tagMap := make(map[string]map[string]interface{})

		for _, analogTag := range deviceConfig.AnalogTagList {
			tagName := analogTag.name.(string)
			if _, exist := a.cfgCache.deviceMap[deviceID][tagName]; !exist && isUpdate {
				continue
			}

			propertyMap := map[string]interface{}{
				"Type": 1,
				"SWVC": analogTag.sendWhenValueChanged,
			}
			if analogTag.spanHigh != nil {
				propertyMap["SH"] = analogTag.spanHigh
			}
			if analogTag.spanLow != nil {
				propertyMap["SL"] = analogTag.spanLow
			}
			if analogTag.fractionDisplayFormat != nil {
				propertyMap["FDF"] = analogTag.fractionDisplayFormat
			}
			if analogTag.integerDisplayFormat != nil {
				propertyMap["IDF"] = analogTag.integerDisplayFormat
			}
			tagMap[tagName] = propertyMap
		}

		for _, discreteTag := range deviceConfig.DiscreteTagList {
			tagName := discreteTag.name.(string)
			if _, exist := a.cfgCache.deviceMap[deviceID][tagName]; !exist && isUpdate {
				continue
			}

			propertyMap := map[string]interface{}{
				"Type": 2,
				"SWVC": discreteTag.sendWhenValueChanged,
			}
			tagMap[tagName] = propertyMap
		}

		for _, textTag := range deviceConfig.TextTagList {
			tagName := textTag.name.(string)
			if _, exist := a.cfgCache.deviceMap[deviceID][tagName]; !exist && isUpdate {
				continue
			}

			propertyMap := map[string]interface{}{
				"Type": 3,
				"SWVC": textTag.sendWhenValueChanged,
			}
			tagMap[tagName] = propertyMap
		}
		newDeviceMap[deviceID] = tagMap
	}

	return newCfgCache
}
