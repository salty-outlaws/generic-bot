# Generic Bot (generic-bot)

A generic bot which uses plugins to add functionality.

The bot makes use of a plugins folder to store lua based plugins. These plugins are loaded at runtime to add functionality to the bot. Each plugin is expected to have a `RegisterCommands` function which registers the commands and functions of the specific plugin to the bot.

## Running locally
```
go run main.go
```

## TODO
- Create basic plugin system
- Create a basic discord bot
- Integrate Discord bot inputs to plugin commands
- Allow plugins to use json, resty, stdlibs and handle db calls
- Update readme with more info
- Create plugins
- Create examples for plugins

## Planned Plugins
- Animal Pics
- Jokes
- Insult & Compliment commands
- Salty Outlaws github integration
- Monopoly Commands
- link store