local i = 0

function request()
    method = "POST"
    path = "/api/account/register"
    headers = {}
    headers["Content-Type"] = "application/x-www-form-urlencoded"
    body = "email=test_user_" .. tostring(i).. "&password=123456"
    i = i + 1
    return wrk.format(method, path, headers, body)
end