function RegisterCommands()
    RegisterCommand("monopoly", "sell", "Sell")
    RegisterCommand("monopoly", "buy", "Buy")
end

function Sell()
    return "sell"
end

function Buy()
    return "buy"
end
