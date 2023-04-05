package controller

import "github.com/gin-gonic/gin"

const basePath = "api"

func StartServer() error {
	r := gin.Default()
	api := r.Group(basePath)

	a := authController{}
	authGroup := api.Group("auth")
	a.addHandlers(authGroup)

	return r.Run()
}
