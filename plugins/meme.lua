function RegisterCommands(filename)
    register(filename, "yuno", "GetYUNO")
end

function register(filename,command,function_name)
    RegisterCommand(filename, command, function_name)
end

function GetYUNO(username, msg)
    return "https://apimeme.com/meme?meme=Y-U-No&top=&bottom="..stringReplace(msg, " ","+")
end