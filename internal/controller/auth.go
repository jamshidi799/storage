package controller

import "github.com/gin-gonic/gin"

type authController struct {
}

func (a *authController) addHandlers(rg *gin.RouterGroup) {
	rg.POST("/register", a.register)
	rg.POST("/login", a.login)
}

func (a *authController) register(ctx *gin.Context) {

}

func (a *authController) login(ctx *gin.Context) {

}
