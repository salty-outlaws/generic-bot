coll = "mono"

function RegisterCommands(filename)
    RegisterCommand(filename, "sell", "Sell")
    RegisterCommand(filename, "buy", "Buy")
    RegisterCommand(filename, "balance", "Balance")
    RegisterCommand(filename, "gamble", "Gamble")
    RegisterCommand(filename, "beg", "Beg")
end

-- ============
-- utility Functions 
-- ============

math.randomseed(os.time())
function getRandonNumber(startv,endv)
    return math.random(startv, endv)
end

function getUserBalance(username)
    balance = dGet(coll, username.."/balance")
    if balance == "" then
        setUserBalance(username,"100")
        balance = "100"
    end
    return tonumber(balance)
end

function setUserBalance(username, balance)
    dPut(coll, username.."/balance", balance)
end

-- ============
-- commands 
-- ============

function Balance(username, msg)
    return string.format("%s's balance\nbalance: %s",username,getUserBalance(username))
end -- balance

function Beg(username, msg)
    log(os.time())
    lastBeg = dGet(coll, username.."/lastBeg")
    if lastBeg ~= "" and os.difftime(os.time(), tonumber(lastBeg)) < 10 then
        return "You are begging too much. stop it!"
    end
    dPut(coll,username.."/lastBeg", tostring(os.time()))

    begAmount = getRandonNumber(0,200)
    setUserBalance(username, tostring(getUserBalance(username)+begAmount))
    return string.format("%s donated %s to %s, %s","xyz",begAmount,username,"go clean their toilet.")
end -- beg

function Gamble(username, msg)
    msgTable = stringSplit(msg, " ")
    -- get user balance
    balance = getUserBalance(username)
    -- gamble amount 
    amount = 0
    -- see if user entered amount
    amountInput =  #msgTable >= 1 and msgTable[1] or "all"
    if amountInput == "all" then
        amount = tonumber(balance)
    else
        amount = tonumber(amountInput)
        if amount > balance then
            return "You do not have that kind of balance"
        end
    end
    
    if amount < 1 then
        return "You do not have any money to gamble"
    end
    -- lost or won 
    gambleState = ""
    -- shows the calculation 
    gambleResult = ""
    -- balance after the gamble 
    newBalance = 0
    if getRandonNumber(1,100) < 70 then
        -- win 70%
        gambleState = "won"
        winAmount = getRandonNumber(0,amount)
        newBalance = balance + winAmount
        gambleResult = balance.." + "..winAmount.." = "..newBalance
        setUserBalance(username, tostring(newBalance))
    else
        -- lose 30%
        lossAmount = getRandonNumber(0,amount)
        newBalance = balance - lossAmount
        if newBalance < 0 then
            newBalance = 0
            lossAmount = newBalance - balance
        end
        gambleResult = balance.." - "..lossAmount.." = "..newBalance
        setUserBalance(username, tostring(newBalance))
    end

    return string.format("%s %s a gamble\nnew balance: %s", username,gambleState,gambleResult)
end -- gamble

function Sell(username, msg)
    return "sell"
end -- sell

function Buy(username, msg)
    return "buy"
end -- buy
