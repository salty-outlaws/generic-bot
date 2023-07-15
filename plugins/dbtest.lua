function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "db", "DBTest")
end

function DBTest(username, msg)
    data = {wallet = 100, bank = 100}
    mUpsert("dbtest","coll1",username, data)
    return text("test complete")
end