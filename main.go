package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/salty-outlaws/generic-bot/plugin"
	"github.com/salty-outlaws/generic-bot/util"
	log "github.com/sirupsen/logrus"
)

var (
	regCommands = map[string]string{}
	regPrefixes = map[string]string{}
	pm          plugin.PluginManager
)

// getPluginFiles - gets list of plugin files
func getPluginFiles() []string {
	files := []string{}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".lua" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return files
}

// RegisterCommand - register a new command from lua plugin
func RegisterCommand(plugin, prefix, command, function string) {
	log.Infof("registered command %s.%s fn(%s)", plugin, command, function)
	regCommands[command] = function
	regPrefixes[prefix+"/"+command] = plugin
}

func HandleCommand(username string, msg string) (map[string]any, error) {
	msg = strings.TrimSpace(msg)
	msgParts := strings.Split(msg, " ")
	prefix := ""
	command := ""
	if len(msgParts) < 2 {
		command = msgParts[0]
	} else {
		prefix = msgParts[0]
		command = msgParts[1]
	}

	if prefix == "admin" && command == "reload" {
		loadPlugins(pm)
		log.Debug("reloaded plugins")
		return map[string]any{
			"type":    "text",
			"message": "Plugins reloaded",
		}, nil
	}

	if pluginName, ok := regPrefixes[prefix+"/"+command]; ok {
		function, funcOK := regCommands[command]
		if !funcOK {
			return nil, errors.New("function not registered")
		}
		ret := ""
		err := pm.CallUnmarshal(
			&ret,
			pluginName+"."+function,
			username,
			strings.Replace(msg, fmt.Sprintf("%s %s ", prefix, command), "", 1))
		if err != nil {
			log.Debugf("error while handling command %s: %v", command, err)
			return nil, err
		}
		log.Debugf(">%s: %s", command, ret)
		output := map[string]any{}
		json.Unmarshal([]byte(ret), &output)
		return output, nil
	} else {
		log.Debug("Invalid Command")
		return nil, errors.New("invalid command")
	}
}

func AddCommands(pm plugin.PluginManager) {
	pm.SetBulk(map[string]any{
		// rest api calling
		"rGet":    util.RGet,
		"rPut":    util.RPut,
		"rPost":   util.RPost,
		"rPatch":  util.RPatch,
		"rDelete": util.RDelete,

		// db commands
		"dGet":    util.DGet,
		"dPut":    util.DPut,
		"dList":   util.DList,
		"dDelete": util.DDelete,

		"jsonToMap":         util.JsonToMap,
		"jsonListToMapList": util.JsonListToMapList,

		// string utils
		"stringSplit":        util.StringSplit,
		"stringJoin":         util.StringJoin,
		"stringReplaceFirst": util.StringReplaceFirst,
		"stringReplace":      util.StringReplace,

		"tagToId": util.DiscordTagToId,
		"idToTag": util.DiscordIdToTag,

		"random": util.RandomNumber,
		"yesno":  util.YesNo,

		// log from a plugin
		"log": func(msg any) { log.Infof("lua: %v", msg) },

		// output types
		"embed": util.Embed,
		"image": util.Image,
		"text":  util.Text,

		// allow registering commands and prefix from plugin
		"RegisterCommand": RegisterCommand,
	})
}

func loadPlugins(pm plugin.PluginManager) {

	// load all plugin files
	for _, p := range getPluginFiles() {
		_, err := pm.Load("./" + p)
		if err != nil {
			fmt.Println("could not load ", p, err.Error())
		} else {
			fmt.Println("loaded", p)
		}
	}

	log.Debug("Registering plugin commands")
	// call RegisterCommands for each plugin
	pm.Each(func(p plugin.Plugin) {
		pm.Call(p.Name+".RegisterCommands", p.Name)
	})

}

func main() {
	log.SetLevel(log.DebugLevel)
	pm = plugin.NewPluginManager()
	defer pm.Close()

	// Add golang functions as commands to be called from lua
	AddCommands(pm)
	HandleCommand("system", "admin reload")

	webserver()
}

func webserver() {
	r := gin.Default()
	r.POST("/api/message", func(c *gin.Context) {

		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, map[string]string{"error": err.Error()})
		}

		input := map[string]any{}
		err = json.Unmarshal(jsonData, &input)
		if err != nil {
			c.JSON(400, map[string]string{"error": err.Error()})
		}

		username := input["username"].(string)
		msg := input["message"].(string)

		output, err := HandleCommand(username, msg)
		if err != nil {
			log.Error(err)
			c.JSON(400, map[string]string{"error": err.Error()})
		}

		c.JSON(http.StatusOK, output)
	})
	r.Run()
}
