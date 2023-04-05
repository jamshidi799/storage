package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"storage/domain"
)

type controller struct {
	service domain.UserService
}

func NewUserController(rg *gin.RouterGroup, us domain.UserService) {
	handler := &controller{service: us}

	rg.POST("/register", handler.register)
	rg.POST("/login", handler.login)
}

func (c *controller) register(ctx *gin.Context) {
	var req *domain.RegisterRequest
	if err := ctx.BindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := c.service.Register(ctx.Request.Context(), req)
	if err != nil {
		//ctx.JSON(getStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *controller) login(ctx *gin.Context) {
	var req *domain.LoginRequest
	if err := ctx.BindJSON(req); err != nil {
		// todo
		//ctx.JSON(http.StatusBadRequest, domain.ErrBadParamInput)
		return
	}

	res, err := c.service.Login(ctx.Request.Context(), req)
	if err != nil {
		//ctx.JSON(getStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
