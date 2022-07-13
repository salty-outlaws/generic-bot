function RegisterCommands()
    RegisterCommand("animals", "cat", "GetCatPicture")
    RegisterCommand("animals", "dog", "GetDogPicture")
end

function GetCatPicture()
    a = RestGet("https://random.dog/woof.json")
    print(a.String(a))
    return a.String(a)
end

function GetDogPicture()
    return "dog"
end
