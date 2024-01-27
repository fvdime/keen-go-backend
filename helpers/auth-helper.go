package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(ctx *gin.Context, role string) (err error) {
	userType := ctx.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized!!!!")
		return err
	}

	return err
}

func MatchUserTypeToUid(ctx *gin.Context, userId string) (err error) {
	userType := ctx.GetString("user_type")
	uId := ctx.GetString("uid")
	err = nil

	if userType == "USER" && uId != userId {
		err = errors.New("Unauthorized!!!!!!!!!!!")
		return err
	}

	err = CheckUserType(ctx, userType)
	return err
}
