function request()
    method = "GET"
    path = "/api/user/profile"
    headers = {}
    headers["Content-Type"] = "application/x-www-form-urlencoded"
    headers["Cookie"] = 'access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQwMzI4MzQsInVzZXJfaWQiOjg3NjQsImVtYWlsIjoidGVzdF91c2VyXzEifQ.dRyP6Xl4q65yzmwWSjvOexVZ905EFHd1dFffZ24bLKU; Path=/; HttpOnly; Expires=Thu, 25 Apr 2024 08:13:54 GMT;'
    return wrk.format(method, path, headers, body)
end