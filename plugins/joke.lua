function RegisterCommands(filename)
    register(filename, "punchline", "GetRandomPunchline")
end

function register(filename,command,function_name)
    RegisterCommand(filename, command, function_name)
end

function GetRandomPunchline()
    response = jsonToMap(rGet("https://official-joke-api.appspot.com/random_joke"))
    output = response["setup"].."\n"..response["punchline"]
    return output
end