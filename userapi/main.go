package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"loggers"
	"user-api/handler"
)

func main() {
	server := &Server{}
	if err := server.Init(); err != nil {
		panic(err)
	}

	client := handler.NewClient(context.Background(), server.UserinfoClient, logger.NewLogger())

	r := gin.Default()
	r.Use(client.GenRequestId, client.SetTraceData, client.Log)
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
