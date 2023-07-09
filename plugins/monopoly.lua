coll = "mono"

function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "sell", "Sell")
    RegisterCommand(filename, "pls", "buy", "Buy")
    RegisterCommand(filename, "pls", "balance", "Balance")
    RegisterCommand(filename, "pls", "gamble", "Gamble")
    RegisterCommand(filename, "pls", "beg", "Beg")
end

-- ============
-- utility Functions 
-- ============

function getWallet(id)
    wallet = dGet(coll, id.."/wallet")
    if wallet == "" then
        wallet = "100"
        setWallet(id, 100)
    end
    return tonumber(wallet)
end

function setWallet(id, amount)
    dPut(coll, id.."/wallet", tostring(amount))
end

function getBank(id)
    bank = dGet(coll, id.."/bank")
    if bank == "" then
        bank = "100"
        setBank(id, 100)
    end
    return tonumber(bank)
end

function setBank(id, amount)
    dPut(coll, id.."/bank", tostring(amount))
end

function deposit(id, amount)
    wallet = getWallet(id)
    if amount <= wallet then
        wallet = wallet - amount
        bank = getBank(id) + amount
        setWallet(wallet)
        setBank(bank)
    end
end

function withdraw(id, amount)
    bank = getBank(id)
    if amount <= bank then
        bank = bank - amount
        wallet = getWallet(id) + amount
        setWallet(wallet)
        setBank(bank)
    end
end

-- ============
-- commands 
-- ============

function Balance(username, msg)
    fields = {
        ["Wallet"] = tostring(getWallet(username)),
        ["Bank"] = tostring(getBank(username)),
    }

    return embed(
        "Monopoly", 
        "Balance Information "..idToTag(username), 
        fields
    )
end

function Beg(username, msg)
    log(os.time())
    lastBeg = dGet(coll, username.."/lastBeg")
    if lastBeg ~= "" and os.difftime(os.time(), tonumber(lastBeg)) < 10 then
        return embed(
        "Monopoly", 
        string.format("%s, You're begging too much, stop it!", idToTag(username)), 
        {}
    )
    end
    dPut(coll,username.."/lastBeg", tostring(os.time()))

    begAmount = random(0,200)
    setWallet(username, tostring(getWallet(username)+begAmount))

    donated_by = jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_name"))["body"]
    job = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_job"))["body"])

    return embed(
        "Monopoly", 
        string.format("%s donated %s to %s, go %s", donated_by, begAmount, idToTag(username), job), 
        {}
    )
end

function Gamble(username, msg)
    msgTable = stringSplit(msg, " ")
    -- get user balance
    balance = getWallet(username)
    -- gamble amount 
    amount = 0

    line = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_gamble_fail"))["body"])

    -- see if user entered amount
    amountInput =  #msgTable >= 1 and msgTable[1] or "all"
    if amountInput == "all" then
        amount = balance
    else
        amount = tonumber(amountInput)
        if amount > balance then
            return embed(
                "Monopoly", 
                string.format("%s", line), 
                {}
            )
        end
    end
    
    if amount < 1 then
        return "You so broke, you gamble in negative."
    end
    -- lost or won 
    gambleState = ""
    -- shows the calculation 
    gambleResult = ""
    -- balance after the gamble 
    newBalance = 0
    line = ""
    if random(1,100) < 70 then
        -- win 70%
        gambleState = "won"
        line = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_gamble_win"))["body"])
        winAmount = random(1,amount)
        newBalance = balance + winAmount
        gambleResult = balance.." + "..winAmount.." = "..newBalance
        setWallet(username, tostring(newBalance))
    else
        -- lose 30%
        gambleState = "lost"
        line = string.lower(jsonToMap(rGet("https://random-apis-brown.vercel.app/api/random_gamble_loss"))["body"])
        lossAmount = random(1,amount)
        newBalance = balance - lossAmount
        if newBalance < 0 then
            newBalance = 0
            lossAmount = newBalance - balance
        end
        gambleResult = balance.." - "..lossAmount.." = "..newBalance
        setWallet(username, tostring(newBalance))
    end

    return embed(
        "Monopoly", 
        string.format("%s %s a gamble\nnew balance: %s. They %s", idToTag(username), gambleState,gambleResult, line), 
        {}
    )
end -- gamble

function Deposit(username, msg)

end

function Sell(username, msg)
    return "sell"
end -- sell

function Buy(username, msg)
    return "buy"
end -- buy
