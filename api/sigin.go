package api

import (
	"blog_gin/dao"
	"blog_gin/models"
	"blog_gin/service"
	"github.com/gin-gonic/gin"
	"log"
)

func (*Api) Sign(c *gin.Context) {
	//	接收用户账号和密码进行注册
	UserInfo := models.Mima{}
	err := c.ShouldBind(&UserInfo)
	if err != nil {
		log.Println("注册时从页面获取信息失败")
		return
	}
	username := UserInfo.UserName
	userPasswd := UserInfo.Passwd

	user := dao.ExistUser(username)
	if user != nil {
		var res models.Result
		res.Code = -999
		res.Error = "账号已经存在"
		//resultJson, err := json.Marshal(res)
		//c.Set("Content-Type", "application/json")
		//if err != nil {
		//	log.Println("json 序列化失败")
		//	return
		//}
		c.JSON(200, gin.H{
			"code":  res.Code,
			"data":  "",
			"error": res.Error,
		})
		return
	}
	_, err = service.Sign(username, userPasswd)

	if err != nil {
		var res models.Result
		res.Code = -999
		res.Error = err.Error()

		c.JSON(200, gin.H{
			"code":  res.Code,
			"data":  "",
			"error": res.Error,
		})
		return
	}

	loginRes, err := service.Login(username, userPasswd)
	//fmt.Println(loginRes.Usf.Uid, loginRes.Usf.UserName)
	if err != nil {
		var res models.Result
		res.Code = -999
		res.Error = err.Error()

		c.JSON(200, gin.H{
			"code":  res.Code,
			"data":  "",
			"error": res.Error,
		})
		return
	}
	var result models.Result
	result.Code = 200
	result.Error = ""
	result.Data = loginRes

	c.JSON(200, gin.H{
		"code":     result.Code,
		"data":     result.Data,
		"error":    result.Error,
		"userName": loginRes.UserInfo.UserName,
	})

}
