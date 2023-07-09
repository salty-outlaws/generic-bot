function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "dog", "GetDogPicture")
    RegisterCommand(filename, "pls", "cat", "GetCatPicture")
    RegisterCommand(filename, "pls", "random_image", "GetRandomImage")
    RegisterCommand(filename, "pls", "fox", "GetRandomFox")
end

function GetCatPicture()
    return image(
        "https://cataas.com"..jsonToMap(rGet("https://cataas.com/cat?json=true"))["url"])
end

function GetDogPicture()
    return image(jsonToMap(rGet("https://dog.ceo/api/breeds/image/random"))["message"])
end

function GetRandomFox()
    return image(jsonToMap(rGet("https://randomfox.ca/floof"))["image"])
end

function GetRandomImage()
    -- return jsonListToMapList(rGet("https://picsum.photos/v2/list?limit=1"))[1]["download_url"]
    return image("https://picsum.photos/200")
end
