function HandleShortcode(arguments)
    if #arguments == 1 then
        image_src = string.format("/images/data/%s", arguments[1])
        return string.format("![image](%s)", image_src)
    elseif #arguments == 2 then
        image_src = string.format("/images/data/%s", arguments[1])
        return string.format("![%s](%s)", arguments[2], image_src)
    else 
        return ""
    end
end