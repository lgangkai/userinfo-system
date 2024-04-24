local i = 0

function request()
    method = "POST"
    path = "/api/account/login"
    headers = {}
    headers["Content-Type"] = "application/x-www-form-urlencoded"
    headers["Cookie"] = 'access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM5NDIwNzUsInVzZXJfaWQiOjEsImVtYWlsIjoiMTIzQGdtYWlsLmNvbSJ9.NY7_6E3nt_P2X9fTNofscnUPi9nwVeC6WWqZiLiAaz0; Path=/; HttpOnly; Expires=Wed, 24 Apr 2024 07:01:15 GMT;'
    body = "email=test_user_" .. tostring(i).. "&password=123456"
    i = i + 1
    return wrk.format(method, path, headers, body)
end