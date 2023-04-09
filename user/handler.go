package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/domain"
	"strings"
)

type controller struct {
	service domain.UserService
}

func NewUserController(rg *gin.RouterGroup, us domain.UserService) *controller {
	handler := &controller{service: us}

	rg.POST("/register", handler.register)
	rg.POST("/login", handler.login)

	return handler
}

func (c *controller) register(ctx *gin.Context) {
	var req domain.RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *controller) login(ctx *gin.Context) {
	var req domain.LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		// todo
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *controller) JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := extractToken(ctx)
		isValid := c.service.VerifyToken(token)
		if !isValid {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
