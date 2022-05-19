package service

import (
	"blog_gin/dao"
	"blog_gin/models"
	"blog_gin/utils"
	"errors"
	"fmt"
)

func Sign(username, userPasswd string) (*models.LoginRes, error) {
	passwd := utils.Md5Crypt(userPasswd, "mszlu")
	//fmt.Println("********************************")
	//fmt.Println(username, passwd)

	//将新的用户插入到数据库中
	uid := dao.Insertuser(username, passwd)
	if uid == -1 {
		fmt.Println("数据库用户注册插入失败")
		return nil, errors.New("注册用户时出现了错误")
	}

	////生成token
	//token, err := utils.Award(&uid)
	//if err != nil {
	//	log.Println("生成token失败")
	//}
	var userInfo models.UserInfo
	userInfo.Uid = uid
	userInfo.UserName = username
	userInfo.Avatar = "Avatar"

	//var lr = &models.LoginRes{
	//	Token:    token,
	//	UserInfo: userInfo,
	//}

	var lr = &models.LoginRes{
		UserInfo: userInfo,
	}

	return lr, nil
}
