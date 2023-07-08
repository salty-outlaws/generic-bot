function RegisterCommands(filename)
    RegisterCommand(filename, "hello", "Hello")
end

function Hello(username, msg)
    return "Hello "..username.."!"
end

function Hello(username, msg)
    return string.format("Hello <@%s>!", username)
end