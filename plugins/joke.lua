function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "punchline", "GetRandomPunchline")
end

function GetRandomPunchline()
    response = jsonToMap(rGet("https://official-joke-api.appspot.com/random_joke"))
    output = response["setup"].."\n"..response["punchline"]
    return output
end