dbname = "mono"
coll = "users"

function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "balance", "Balance")
    RegisterCommand(filename, "pls", "gamble", "Gamble")
    RegisterCommand(filename, "pls", "beg", "Beg")
    RegisterCommand(filename, "pls", "deposit", "Deposit")
    RegisterCommand(filename, "pls", "draw", "Withdraw")
    RegisterCommand(filename, "pls", "withdraw", "Withdraw")
    RegisterCommand(filename, "pls", "rich", "Rich")
    RegisterCommand(filename, "pls", "search", "Search")
    RegisterCommand(filename, "pls", "steal", "Steal")
    RegisterCommand(filename, "pls", "mug", "Steal")
    RegisterCommand(filename, "pls", "daily", "Daily")
    RegisterCommand(filename, "pls", "send", "Send")
    -- in progress
    RegisterCommand(filename, "pls", "loan", "Loan")
    RegisterCommand(filename, "pls", "bankrob", "BankRob")
    RegisterCommand(filename, "pls", "sell", "Sell")
    RegisterCommand(filename, "pls", "buy", "Buy")
    RegisterCommand(filename, "pls", "shop", "Shop")
    RegisterCommand(filename, "pls", "inventory", "Inventory")
    RegisterCommand(filename, "pls", "use", "Use")
    RegisterCommand(filename, "pls", "lottery", "Lottery")
end

-- ============
-- utility Functions 
-- ============
function dbSet(id, data)
    return mUpsert(dbname, coll, id, data)
end

function dbGet(id)
    return mGet(dbname, coll, id)
end

function dbDelete(id)
    return mDelete(dbname, coll, id)
end

function dbFind(query, sort)
    return mFind(dbname, coll, query, sort)
end

function createUser(id)
    data = {wallet = 100, bank = 100}
    dbSet(id, data)
end

function hasUser(id)
    return dbGet(id) ~= nil
end

function getWallet(id)
    if hasUser(id) == false then
        createUser(id)
    end
    return dbGet(id)["wallet"]
end

function setWallet(id, amount)
    if hasUser(id) == false then
        createUser(id)
        amount = amount + 100
    end
    u = dbGet(id)
    u["wallet"] = amount
    dbSet(id, u)
end

function getBank(id)
    if hasUser(id) == false then
        createUser(id)
    end
    return dbGet(id)["bank"]
end

function setBank(id, amount)
    if hasUser(id) == false then
        createUser(id)
        amount = amount + 100
    end
    u = dbGet(id)
    u["bank"] = amount
    dbSet(id, u)
end

function depositWallet(id, amount)
    wallet = getWallet(id)
    if amount <= wallet then
        wallet = wallet - amount
        bank = getBank(id) + amount
        setWallet(wallet)
        setBank(bank)
    end
end

function withdrawBank(id, amount)
    bank = getBank(id)
    if amount <= bank then
        bank = bank - amount
        wallet = getWallet(id) + amount
        setWallet(wallet)
        setBank(bank)
    end
end

-- embed partial function
function emb(message)
    return embed(
        "Monopoly",
        message,
        {}
    )
end

function parseAmount(id, msg, balance, index)
    msgTable = stringSplit(msg, " ")
    amount = 0

    if msg == "" then
        amount = balance
    elseif #msgTable >= index then
        if msgTable[index] == "all" then
            amount = balance
        else
            amount = tonumber(msgTable[index])
        end
    end
    return amount
end
-- ============
-- commands 
-- ============

function Balance(id, msg)
    fields = {
        ["Wallet"] = tostring(getWallet(id)),
        ["Bank"] = tostring(getBank(id)),
    }

    return embed( 
        "Monopoly",
        "Balance Information "..idToTag(id), 
        fields
    )
end

function Beg(id, msg)
    -- beg in range 1-100
    u = dbGet(id)
    if u["last_beg"] ~= nil then
        if u["last_beg"] ~= "" and os.difftime(os.time(), u["last_beg"]) < 10 then
            return emb(string.format("%s, You're begging too much, stop it!", idToTag(id)))
        end
    end
    u["last_beg"] = os.time()
    dbSet(id, u)

    begAmount = random(1,200)
    donated_by = jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_name"))["body"]
    job = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_job"))["body"])
    
    setWallet(id, getWallet(id)+begAmount)
    
    return emb(string.format("%s donated $%d to %s, go %s", donated_by, begAmount, idToTag(id), job))
end

function Gamble(id, msg)
    balance = getWallet(id)
    amount = parseAmount(id, msg, balance, 1)

    if amount < 1 then
        return emb("We are here to gamble, not play with imaginary numbers!")
    end

    -- cant gamble if you dont have money in wallet 
    gamble_fail_line = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_gamble_fail"))["body"])
    if amount > balance then
        return emb(string.format("%s", gamble_fail_line))
    end

    -- lost or won
    gambleState = ""
    gambleResult = ""
    newBalance = 0
    line = ""
    if random(1,100) < 60 then
        -- win 60%
        gambleState = "won"
        line = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_gamble_win"))["body"])
        winAmount = random(1,amount)
        newBalance = balance + winAmount
        gambleResult = balance.." + "..winAmount.." = "..newBalance
        setWallet(id, newBalance)
    else
        -- lose 40%
        gambleState = "lost"
        line = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_gamble_loss"))["body"])
        lossAmount = random(1,amount)
        newBalance = balance - lossAmount
        if newBalance < 0 then
            newBalance = 0
            lossAmount = newBalance - balance
        end
        gambleResult = balance.." - "..lossAmount.." = "..newBalance
        setWallet(id, newBalance)
    end

    return emb(string.format("%s %s a gamble\nnew balance: %s. They %s", idToTag(id), gambleState,gambleResult, line))
end

function Deposit(id, msg)
    balance = getWallet(id)
    amount = parseAmount(id, msg, balance, 1)
    if amount < 1 or amount > balance then
        emb("Cannot deposit that amount.")
    end
    
    setBank(id, getBank(id)+amount)
    setWallet(id, balance-amount)
    return emb(string.format("Deposited %s in %s's bank",amount,idToTag(id)))
end

function Withdraw(id, msg)
    balance = getBank(id)
    amount = parseAmount(id, msg, balance, 1)
    if amount < 1 or amount > balance then
        emb("Cannot withdraw that amount.")
    end
    
    setWallet(id, getWallet(id)+amount)
    setBank(id, balance-amount)
    return emb(string.format("Withdrew %s from %s's bank",amount,idToTag(id)))
end

function Rich(id, msg)
    values = dbFind({}, {wallet = -1})
    output = "Rich List"
    for i = 1, 10 do
        if #values < i then
            break
        end
        output = string.format("%s\n%s = $%d", output, idToTag(values[i]["_id"]), values[i]["wallet"])
    end
    return emb(output)
end

function Search(id, msg)
    -- search with less chances 200-300
    u = dbGet(id)
    if u["last_search"] ~= nil then
        if u["last_search"] ~= "" and os.difftime(os.time(), u["last_search"]) < 30 then
            return emb(string.format("%s, You're tired of searching, take some rest!", idToTag(id)))
        end
    end
    u["last_search"] = os.time()
    dbSet(id, u)

    searchAmount = random(200,300)
    place = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_money_place"))["body"])
    
    setWallet(id, getWallet(id)+searchAmount)
    return emb(string.format("%s found $%d %s", idToTag(id), searchAmount, place))
end

function Sell(id, msg)
    return "sell"
end -- sell

function Buy(id, msg)
    return "buy"
end -- buy

function Steal(id, msg)
    msgTable = stringSplit(msg, " ")

    -- get user and victim ids
    user = id
    victim = tagToId(msgTable[1])

    -- get user and victim db objects
    u = dbGet(user)
    v = dbGet(victim)

    -- ensure both arent broke
    if u["wallet"] < 200 or v["wallet"] < 200 then
        return emb("You need atleast 200 in victim's and your wallet")
    end

    -- establish respectable amount to steal 
    stealCap = u["wallet"]
    if stealCap > v["wallet"] then
        stealCap = v["wallet"]
    end
    stealCap = stealCap/2
    stealAmount = random(1, stealCap)
    
    -- attempt to steal 
    result = ""
    if yesno() == true then
        mugging_location = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_mug_place"))["body"])
        v["wallet"] = v["wallet"] - stealAmount
        u["wallet"] = u["wallet"] + stealAmount    
        result = string.format("%s stole $%d from %s %s", idToTag(user), stealAmount, idToTag(victim), mugging_location)
    else
        u["wallet"] = u["wallet"] - stealAmount
        v["wallet"] = v["wallet"] + stealAmount
        result = string.format("%s got caught while stealing from %s and paid them $%d", idToTag(user), idToTag(victim), stealAmount)
    end
    setWallet(user, u["wallet"])
    setWallet(victim, v["wallet"])

    return emb(result)
end

function Shop(id, msg)
    return "Shop"
end

function Inventory(id, msg)
    return "Inventory"
end

function Use(id, msg)
    return "Use"
end

function Send(id, msg)
    msgTable = stringSplit(msg, " ")

    -- get user and victim ids
    sender = id
    receiver = tagToId(msgTable[1])

    s = dbGet(sender)
    r = dbGet(receiver)

    amount = tonumber(msgTable[2])
    if amount > s["wallet"] then
        return emb("You're too broke for that.")
    end

    s["wallet"] = s["wallet"] - amount
    r["wallet"] = r["wallet"] + amount

    setWallet(sender, s["wallet"])
    setWallet(receiver, r["wallet"])

    return emb(string.format("%s sent $%d to %s", idToTag(sender), amount, idToTag(receiver)))
end

function Loan(id, msg)
    return "Loan"
end

function Daily(id, msg)
    u = dbGet(id)
    if u["last_daily"] ~= nil then
        if u["last_daily"] ~= "" and os.difftime(os.time(), u["last_daily"]) < 86400 then
            return emb(string.format("%s, You already got your daily allowance!", idToTag(id)))
        end
    end
    u["last_daily"] = os.time()
    dbSet(id, u)

    dailyAmount = random(1000,1500)    
    setWallet(id, getWallet(id)+dailyAmount)
    return emb(string.format("%s got a daily allowance of $%d", idToTag(id), dailyAmount))
end

function Lottery(id, msg)
    return "Lottery"
end

function BankRob(id, msg)
    return "BankRob"
end