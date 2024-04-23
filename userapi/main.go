package main

import (
	"github.com/gin-gonic/gin"
	"loggers"
	"user-api/handler"
)

func main() {
	log := logger.NewLogger()
	context := &gin.Context{}
	context.Set("request_id", 123456)
	context.Set("user_id", 123)
	log.Debug(context, "test", 12345667890)
	server := &Server{}
	if err := server.Init(); err != nil {
		panic(err)
	}

	client := &handler.Client{UserinfoClient: server.UserinfoClient}

	r := gin.Default()
	r.Use(handler.Log)
	apiAccount := r.Group("api/account")
	{
		apiAccount.POST("login", client.Login)
		apiAccount.POST("logout", client.Logout)
		apiAccount.POST("register", client.Register)
	}

	apiProfile := r.Group("api/user/profile")
	apiProfile.Use(client.Authenticate)
	{
		apiProfile.GET("", client.GetProfile)
		apiProfile.POST("", client.CreateProfile)
		apiProfile.PUT("", client.UpdateProfile)
	}

	if err := r.Run(server.Addr); err != nil {
		panic(err)
	}
}
