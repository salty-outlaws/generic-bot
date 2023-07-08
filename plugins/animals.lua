function RegisterCommands(filename)
    register(filename, "dog", "GetDogPicture")
    register(filename, "cat", "GetCatPicture")
    register(filename, "random_image", "GetRandomImage")
end

function register(filename,command,function_name)
    RegisterCommand(filename, command, function_name)
end

function GetCatPicture()
    return "https://cataas.com"..jsonToMap(rGet("https://cataas.com/cat?json=true"))["url"]
end

function GetDogPicture()
    return jsonToMap(rGet("https://dog.ceo/api/breeds/image/random"))["message"]
end

function GetRandomImage()
    -- return jsonListToMapList(rGet("https://picsum.photos/v2/list?limit=1"))[1]["download_url"]
    return "https://picsum.photos/200"
end