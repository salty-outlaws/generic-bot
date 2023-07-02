function RegisterCommands(filename)
    register(filename, "dog", "GetDogPicture")
    register(filename, "cat", "GetCatPicture")
end

function register(filename,command,function_name)
    RegisterCommand(filename, command, function_name)
end

function GetCatPicture()
    -- log("getting cat pics")
    return rGet("https://cataas.com/cat?json=true")
end

function GetDogPicture()
    return rGet("https://dog.ceo/api/breeds/image/random")
end
