package logic

import (
	"Bluebell/dao/mysql"
	"Bluebell/models"
	"Bluebell/pkg/jwt"
	"Bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户存不存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}
	//生成UID
	userID := snowflake.GenID()
	//构造User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就能拿到User.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT的token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
