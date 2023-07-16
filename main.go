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
	"github.com/salty-outlaws/generic-bot/plugin"
	"github.com/salty-outlaws/generic-bot/util"
	log "github.com/sirupsen/logrus"
)

var (
	pluginDirectory = map[string]string{}
	pm              plugin.PluginManager
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
	log.Infof("registered command %s/%s -> %s.%s", prefix, command, plugin, function)
	pluginDirectory[getPluginKey(prefix, command)] = strings.Join([]string{plugin, function}, ".")
}

func getPluginKey(prefix, command string) string {
	return strings.Join([]string{prefix, command}, "/")
}

func HandleCommand(guild, username, msg string) (map[string]any, error) {
	msgs := strings.Split(strings.TrimSpace(msg), " ")
	prefix := ""
	command := ""
	if len(msgs) == 1 {
		command = msgs[0]
	} else {
		prefix = msgs[0]
		command = msgs[1]
	}

	if prefix == "admin" && command == "reload" {
		loadPlugins()
		return map[string]any{
			"type":    "text",
			"message": "Plugins reloaded",
		}, nil
	}

	pluginEntry, ok := pluginDirectory[getPluginKey(prefix, command)]
	if !ok {
		errLine := fmt.Sprintf("No plugin entry found for prefix %s command %s", prefix, command)
		log.Debug(errLine)
		return map[string]any{"type": "debug", "message": errLine}, nil
	}

	msg = strings.Replace(msg, prefix, "", 1)
	msg = strings.Replace(msg, command, "", 1)
	msg = strings.TrimSpace(msg)

	ret := ""
	err := pm.CallUnmarshal(
		&ret, pluginEntry, username, msg)
	if err != nil {
		log.Errorf("error while handling command %s: %v", command, err)
		return nil, err
	}

	log.Debugf(">%s: %s", getPluginKey(prefix, command), ret)
	output := map[string]any{}
	json.Unmarshal([]byte(ret), &output)
	return output, nil
}

func AddCommands(pm plugin.PluginManager) {
	pm.SetBulk(map[string]any{
		// rest api calling
		"rGet":    util.RGet,
		"rPut":    util.RPut,
		"rPost":   util.RPost,
		"rPatch":  util.RPatch,
		"rDelete": util.RDelete,

		"mUpsert": util.MUpsert,
		"mGet":    util.MGet,
		"mDelete": util.MDelete,
		"mFind":   util.MFind,

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

func loadPlugins() {

	resetPm()

	// load all plugin files
	for _, p := range getPluginFiles() {
		_, err := pm.LoadFile("./" + p)
		if err != nil {
			log.Errorf("could not load %s error %s", p, err.Error())
		}
	}

	pluginReposFile, _ := ioutil.ReadFile("plugin_repos.json")
	pluginRepos := []string{}
	err := json.Unmarshal(pluginReposFile, &pluginRepos)
	if err != nil {
		log.Panicf("Could not load repos file")
	}

	for _, url := range pluginRepos {
		repoPlugins, err := plugin.LoadPluginRepo(url)
		if err != nil {
			log.Errorf("Could not load plugins from %s", url)
		}
		for _, pluginName := range repoPlugins {
			_, err := pm.LoadUrl(fmt.Sprintf("%s/%s?raw=true", url, pluginName))
			if err != nil {
				log.Errorf("could not load URL: %s", err.Error())
			}
		}
	}

	log.Debug("Registering plugin commands")
	// call RegisterCommands for each plugin
	pm.Each(func(p plugin.Plugin) {
		pm.Call(p.Name+".RegisterCommands", p.Name)
	})

}

func resetPm() {
	pm = plugin.NewPluginManager()
	AddCommands(pm)
}

func main() {
	log.SetLevel(log.DebugLevel)
	// Add golang functions as commands to be called from lua
	HandleCommand("", "system", "admin reload")

	webserver()
}

func webserver() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/api/message", func(c *gin.Context) {

		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, map[string]string{"error": err.Error()})
			return
		}

		input := map[string]any{}
		err = json.Unmarshal(jsonData, &input)
		if err != nil {
			c.JSON(400, map[string]string{"error": err.Error()})
			return
		}

		guild := input["guild"].(string)
		username := input["username"].(string)
		msg := input["message"].(string)

		output, err := HandleCommand(guild, username, msg)
		if err != nil {
			log.Error(err)
			c.JSON(400, map[string]string{
				"type":  "error",
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, output)
	})
	r.Run()
}
