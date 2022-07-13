package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/salty-outlaws/generic-bot/plugin"
)

var (
	restClient  *resty.Client
	regCommands = map[string]RegisteredCommandParams{}
)

// RegisteredCommandParams - params for each command
type RegisteredCommandParams struct {
	Plugin   string
	Function string
}

// getPlugins - send list of plugin files
func getPlugins() []string {
	files := []string{}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".lua" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return files
}

func RegisterCommand(plugin string, command string, function string) {
	fmt.Println("registered", plugin, command, function)
	regCommands[command] = RegisteredCommandParams{
		Plugin:   plugin,
		Function: function,
	}
}

func HandleCommand(pm plugin.PluginManager, msg string) {
	command := strings.Split(msg, " ")[0]
	plugParams := regCommands[command]
	// ret, err := pm.Call(plugParams.Plugin + "." + plugParams.Function)
	ret := ""
	err := pm.CallUnmarshal(&ret, plugParams.Plugin+"."+plugParams.Function)
	_ = err
	fmt.Println(ret, err)
}

func ImportResty(pm plugin.PluginManager) {
	r := restClient.R()
	pm.SetBulk(map[string]any{
		"RestGet":    r.Get,
		"RestPut":    r.Put,
		"RestPost":   r.Post,
		"RestDelete": r.Delete,
	})
}

func main() {

	restClient = resty.New()
	pm := plugin.NewPluginManager()
	defer pm.Close()

	ImportResty(pm)

	for _, p := range getPlugins() {
		_, err := pm.Load("./" + p)
		if err != nil {
			fmt.Println("could not load ", p, err.Error())
		} else {
			fmt.Println("loaded", p)
		}
	}
	pm.Set("RegisterCommand", RegisterCommand)

	pm.Each(func(p plugin.Plugin) {
		pm.Call(p.Name + ".RegisterCommands")
	})

	HandleCommand(pm, "cat pics")
	HandleCommand(pm, "dog pics")
	HandleCommand(pm, "sell soul")
}
