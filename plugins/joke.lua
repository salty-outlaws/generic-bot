function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "punchline", "GetRandomPunchline")
    RegisterCommand(filename, "pls", "joke", "GetRandomPunchline")
end

function GetRandomPunchline()
    response = jsonToMap(rGet("https://official-joke-api.appspot.com/random_joke"))
    output = response["setup"].."\n"..response["punchline"]
    return text(output)
end