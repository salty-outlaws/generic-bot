package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/salty-outlaws/generic-bot/util"
)

func LoadPluginRepo(url string) ([]string, error) {
	configJson := util.RGet(fmt.Sprintf("%s/config.json?raw=true", url))
	configMap := []string{}
	err := json.Unmarshal([]byte(configJson), &configMap)
	return configMap, err
}
