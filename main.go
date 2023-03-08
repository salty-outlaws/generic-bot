package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/salty-outlaws/generic-bot/db"
	"github.com/salty-outlaws/generic-bot/plugin"
	"github.com/salty-outlaws/generic-bot/rest"
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

func HandleCommand(username string, msg string) {
	command := strings.Split(msg, " ")[0]
	plugParams := regCommands[command]
	// ret, err := pm.Call(plugParams.Plugin + "." + plugParams.Function)
	ret := ""
	err := pm.CallUnmarshal(&ret, plugParams.Plugin+"."+plugParams.Function, username, strings.Split(msg, " "))
	if err != nil {
		log.Errorf("error while handling command %s: %v", command, err)
		return
	}
	log.Debugf(">%s: %s", command, ret)
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

		// log from a plugin
		"log": func(msg any) { log.Infof("lua: %v", msg) },

		// allow registering commands from plugin
		"RegisterCommand": RegisterCommand,
	})
}

func main() {
	log.SetLevel(log.DebugLevel)

	pm = plugin.NewPluginManager()
	defer pm.Close()

	// load all plugin files
	for _, p := range getPluginFiles() {
		_, err := pm.Load("./" + p)
		if err != nil {
			fmt.Println("could not load ", p, err.Error())
		} else {
			fmt.Println("loaded", p)
		}
	}

	// Add golang functions as commands to be called from lua
	AddCommands(pm)

	// call RegisterCommands for each plugin
	pm.Each(func(p plugin.Plugin) {
		pm.Call(p.Name+".RegisterCommands", p.Name)
	})

	// test code
	test()
}

func test() {
	for {
		var input string
		fmt.Scanln(&input)
		HandleCommand("mridulganga", input)
	}

	// HandleCommand("mridulganga", "cat pics")
	// HandleCommand("mridulganga", "dog pics")
	// HandleCommand("mridulganga", "balance")
	// // HandleCommand("mridulganga", "gamble all")
	// // HandleCommand("mridulganga", "gamble all")
	// // HandleCommand("mridulganga", "gamble all")
	// // HandleCommand("mridulganga", "gamble all")
	// HandleCommand("mridulganga", "beg")
}
