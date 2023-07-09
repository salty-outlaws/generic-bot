function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "ping", "Ping")
end

function Ping(username, msg)
    log("pinging "..msg)
    local handler = io.popen("ping -c 3 -i 0.5 "..msg)
    local response = handler:read("*a")
    return text(response)
end