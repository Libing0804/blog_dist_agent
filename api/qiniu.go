package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func (*Api) QiniuToken(c *gin.Context) {
	//	自定义凭证有效时间
	bucket := "blogpicturesave"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	putPolicy.Expires = 7200 //两小时
	//mac := qbox.NewMac(config.Cfg.System.QiniuAccessKey, config.Cfg.System.QiniuSecretKey)
	mac := qbox.NewMac("OiS_p1_Md9ShK8y9kkMWp2jzSVDpMv--Cr5IT-ty", "X9QchF4lb4Pe7MDFaVTvePGMLnuwlgcY_r0tdQNs")

	upToken := putPolicy.UploadToken(mac)
	c.JSON(200, gin.H{
		"code":  200,
		"error": "",
		"data":  upToken,
	})
}
