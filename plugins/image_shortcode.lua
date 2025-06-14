function HandleShortcode(arguments)
    if #arguments == 1 then
        image_src = string.format("/image/%s", arguments[1])
        return string.format("![image](/image/%s)", arguments[1])
    elseif #arguments == 2 then
        image_src = string.format("/image/%s", arguments[1])
        return string.format("![%s](/image/%s)", arguments[2], arguments[1])
    else 
        return ""
    end
end