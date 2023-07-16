function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "insult", "GetRandomInsult")
    RegisterCommand(filename, "pls", "compliment", "GetRandomCompliment")
end

function GetRandomInsult(id, msg)
    msgTable = stringSplit(msg, " ")

    -- Get ID of person to insult
    -- Insult the sender unless they have tagged someone
    victim = id
    if msgTable[1] ~= nil then
        victim = tagToId(msgTable[1])
    end

    insult = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_insult"))["body"])
    output = string.format("Hey %s, %s", idToTag(victim), insult)
    return text(output)
end

function GetRandomCompliment(id, msg)
    msgTable = stringSplit(msg, " ")

    -- Get ID of person to insult
    -- Insult the sender unless they have tagged someone
    victim = id
    if msgTable[1] ~= nil then
        victim = tagToId(msgTable[1])
    end

    insult = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_compliment"))["body"])
    output = string.format("Hey %s, %s", idToTag(victim), insult)
    return text(output)
end