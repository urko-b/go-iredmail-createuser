package middleware

import (
	"fmt"
	"iredmail-create-email-account/pkg/public_error"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Error string `json:"error"`
}

func newError(err error) ApiError {
	return ApiError{
		Error: err.Error(),
	}
}

func Error(support_email string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		responseErr := newError(
			fmt.Errorf(
				"couldn't create user, try again or contact administrator: %s", support_email,
			),
		)

		for i, ginErr := range ctx.Errors {
			fmt.Printf(
				"error: %d\n json gin error: %s\n",
				i,
				ginErr.Error(),
			)
			fmt.Println(fmt.Sprintf("ginErr %d", i), ginErr)

			if public_error.IsPublicErr(ginErr.Err) {
				publicErr, _ := ginErr.Err.(public_error.Public)
				responseErr = newError(publicErr)
			}
		}
		ctx.JSON(
			-1,
			responseErr,
		)
	}
}
