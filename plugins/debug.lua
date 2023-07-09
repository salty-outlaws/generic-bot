function RegisterCommands(filename)
    RegisterCommand(filename, "pls","log", "LogMsg")
end

function LogMsg(username, msg)
    return msg
end