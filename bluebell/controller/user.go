package controller

import (
	"Bluebell/dao/mysql"
	"Bluebell/logic"
	"Bluebell/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	//"github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid", zap.Error(err))
		//判断err是不是vikidator.ValidatorErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	fmt.Println(p)

	//业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应

	ResponseSuccess(c, nil)

}

func LoginHandler(c *gin.Context) {
	//处理请求参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindBodyWithJSON(p); err != nil {
		zap.L().Error("Login with invalid", zap.Error(err))
		//判断err是不是vikidator.ValidatorErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}
	//业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	//返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
