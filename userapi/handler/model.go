package handler

type Profile struct {
	Username string `form:"username"`
	Birthday string `form:"birthday"`
}

type Account struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
