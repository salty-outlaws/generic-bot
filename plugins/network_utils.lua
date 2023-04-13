function RegisterCommands(filename)
    register(filename, "ping", "Ping")
end

function register(filename,command,function_name)
    RegisterCommand(filename, command, function_name)
end

function Ping(username, msg)
    log("pinging "..msg[2])
    local handler = io.popen("ping -c 3 -i 0.5 "..msg[2])
    local response = handler:read("*a")
    return response
end