function RegisterCommands(filename)
    register(filename, "dog", "GetDogPicture")
    register(filename, "cat", "GetCatPicture")
end

function register(filename,command,function_name)
    RegisterCommand(filename, command, function_name)
end

function GetCatPicture()
    -- log("getting cat pics")
    return rGet("https://random.dog/woof.json")
end

function GetDogPicture()
    return rGet("https://dog.ceo/api/breeds/image/random")
end
