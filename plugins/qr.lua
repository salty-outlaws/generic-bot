function RegisterCommands(filename)
    RegisterCommand(filename, "pls", "qr", "QRCodeGenerate")
end

function QRCodeGenerate(username, msg)
    return image("https://api.qrserver.com/v1/create-qr-code/?data="..stringReplace(msg, " ", "%20"))
end