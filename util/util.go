package util

import (
	"encoding/json"
	"strings"
)

func JsonToMap(val string) map[string]any {
	output := map[string]any{}
	json.Unmarshal([]byte(val), &output)
	return output
}

func JsonListToMapList(val string) []map[string]any {
	output := []map[string]any{}
	json.Unmarshal([]byte(val), &output)
	return output
}

func StringSplit(value string, delimeter string) []string {
	return strings.Split(value, delimeter)
}

func StringReplaceFirst(value, old, new string) string {
	return strings.Replace(value, old, new, 1)
}

func StringReplace(value, old, new string) string {
	return strings.ReplaceAll(value, old, new)
}
