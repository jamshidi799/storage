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
	var req registerRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.service.Register(ctx.Request.Context(), req.toUser())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, registerResponse{
		Id:    user.Id,
		Email: user.Email,
	})
}

func (c *controller) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.BindJSON(&req); err != nil {
		// todo
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.service.Login(ctx.Request.Context(), req.toUser())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, loginResponse{Token: token})
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

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *registerRequest) toUser() *domain.User {
	return &domain.User{
		Email:    r.Email,
		Password: r.Password,
	}
}

type registerResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (l *loginRequest) toUser() *domain.User {
	return &domain.User{
		Email:    l.Email,
		Password: l.Password,
	}
}

type loginResponse struct {
	Token string `json:"token"`
}
