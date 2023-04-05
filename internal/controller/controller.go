package controller

import "github.com/gin-gonic/gin"

func StartServer() error {
	r := gin.Default()

	return r.Run()
}
