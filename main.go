package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/salty-outlaws/generic-bot/db"
	"github.com/salty-outlaws/generic-bot/plugin"
	"github.com/salty-outlaws/generic-bot/rest"
	"github.com/salty-outlaws/generic-bot/util"
	log "github.com/sirupsen/logrus"
)

var (
	regCommands = map[string]RegisteredCommandParams{}
	pm          plugin.PluginManager
)

// RegisteredCommandParams - params for each command
type RegisteredCommandParams struct {
	Plugin   string
	Function string
}

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
func RegisterCommand(plugin string, command string, function string) {
	log.Infof("registered command %s.%s fn(%s)", plugin, command, function)
	regCommands[command] = RegisteredCommandParams{
		Plugin:   plugin,
		Function: function,
	}
}

func HandleCommand(username string, msg string) (string, error) {
	msg = strings.TrimSpace(msg)
	command := strings.Split(msg, " ")[0]

	switch command {
	case "":
		log.Debug("Invalid Command")
	case "reload":
		loadPlugins(pm)
		log.Debug("reloaded plugins")
	default:
		plugParams := regCommands[command]
		// ret, err := pm.Call(plugParams.Plugin + "." + plugParams.Function)
		ret := ""
		err := pm.CallUnmarshal(
			&ret,
			plugParams.Plugin+"."+plugParams.Function,
			username,
			strings.Replace(msg, command+" ", "", 1))
		if err != nil {
			log.Errorf("error while handling command %s: %v", command, err)
			return "", err
		}
		log.Debugf(">%s: %s", command, ret)
		return ret, nil
	}
	return "", nil
}

func AddCommands(pm plugin.PluginManager) {
	pm.SetBulk(map[string]any{
		// rest api calling
		"rGet":    rest.Get,
		"rPut":    rest.Put,
		"rPost":   rest.Post,
		"rDelete": rest.Delete,

		// db commands
		"dGet":    db.Get,
		"dPut":    db.Put,
		"dList":   db.List,
		"dDelete": db.Delete,

		"jsonToMap":         util.JsonToMap,
		"jsonListToMapList": util.JsonListToMapList,

		"stringSplit":        util.StringSplit,
		"stringReplaceFirst": util.StringReplaceFirst,
		"stringReplace":      util.StringReplace,

		// log from a plugin
		"log": func(msg any) { log.Infof("lua: %v", msg) },

		// allow registering commands from plugin
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
	HandleCommand("system", "reload")

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
			c.JSON(400, map[string]string{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": output,
		})
	})
	r.Run()
}
