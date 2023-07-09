function RegisterCommands(filename)
    RegisterCommand(filename, "pls","yuno", "GetYUNO")
    RegisterCommand(filename, "pls","zuck", "GetZuck")
end

function GetYUNO(username, msg)
    return image("https://apimeme.com/meme?meme=Y-U-No&top=Y+U+NO&bottom="..stringReplace(msg, " ","+"))
end

function GetZuck(username, msg)
    msg = stringReplace(msg, " ", "+")
    top = stringSplit(msg, ":")[1]
    bottom = stringSplit(msg, ":")[2]
    return image("https://apimeme.com/meme?meme=Zuckerberg&top="..top.."&bottom="..bottom)
end