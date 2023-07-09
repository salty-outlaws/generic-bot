function RegisterCommands(filename)
    RegisterCommand(filename, "", "hello", "Hello")
end

function Hello(username, msg)
    return text(string.format("Hello <@%s>!", username))
end