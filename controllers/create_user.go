package controllers

import (
	"iredmail-create-email-account/pkg/create_user"
	"iredmail-create-email-account/pkg/public_error"
	"iredmail-create-email-account/pkg/remote_ssh"
	"iredmail-create-email-account/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	create_user_service create_user.CreateUserService
}

func NewUser(config remote_ssh.Config) UserController {
	utils.LogJSONObject("remote_config", config)
	return UserController{
		create_user_service: create_user.New(
			config.KeyPath,
			config.Server,
			config.User,
			config.Port,
		),
	}
}

func (ctrl UserController) CreateUser(ctx *gin.Context) {
	dto := &create_user.CreateUserDTO{}
	err := ctx.Bind(dto)
	if err != nil {
		wrapAbortWithError(ctx, err)
		return
	}

	err = ctrl.create_user_service.Create(dto)
	if err != nil {
		if public_error.IsPublicErr(err) {
			publicErr, _ := err.(public_error.Public)
			wrapAbortStatusCodeWithError(ctx, publicErr.StatusCode(), publicErr)
			return
		}
		wrapAbortWithError(ctx, err)
	}
}

func wrapAbortWithError(ctx *gin.Context, err error) {
	ctx.AbortWithError(http.StatusInternalServerError, err) //nolint: errcheck
}

func wrapAbortStatusCodeWithError(ctx *gin.Context, status_code int, err error) {
	ctx.AbortWithError(status_code, err) //nolint: errcheck
}
