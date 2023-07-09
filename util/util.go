package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
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

func StringJoin(values []string, delimeter string) string {
	return strings.Join(values, delimeter)
}

func StringReplaceFirst(value, old, new string) string {
	return strings.Replace(value, old, new, 1)
}

func StringReplace(value, old, new string) string {
	return strings.ReplaceAll(value, old, new)
}

/*
Converts Discord Tag to Discord Id
ex - <@1011508634460631100> to 1011508634460631100
*/
func DiscordTagToId(value string) string {
	output := strings.Replace(value, "<@", "", 1)
	output = strings.Replace(output, ">", "", 1)
	return output
}

/*
Converts Discord Id to Discord Tag
ex - 1011508634460631100 to <@1011508634460631100>
*/
func DiscordIdToTag(value string) string {
	return fmt.Sprintf("<@%s>", value)
}

/*
Generates random integer between from and to
*/
func RandomNumber(from int, to int) int {
	rand.Seed(time.Now().UnixNano())
	return from + rand.Intn(to)
}

func YesNo() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
