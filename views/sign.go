package views

import (
	"blog_gin/config"
	"github.com/gin-gonic/gin"
)

func (*HTMLApi) Sign(c *gin.Context) {
	c.HTML(200, "sign.html", config.Cfg.Viewer)
}
