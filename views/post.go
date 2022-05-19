package views

import (
	"blog_gin/dao"
	"blog_gin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (*HTMLApi) About(c *gin.Context) {
	PostRes, err := service.GetDetailPost(102)
	if err != nil {
		c.HTML(http.StatusOK, "detail.html", "文章查询失败")
		return
	}
	//每一次查看文章都应该将查看次数加1
	dao.UpdateViewCount(102, PostRes.Article.ViewCount+1)
	c.HTML(http.StatusOK, "detail.html", PostRes)
}
func (*HTMLApi) Detail(c *gin.Context) {
	pidstr := c.Param("pid")
	pidstr = strings.TrimSuffix(pidstr, ".html")
	pid, err := strconv.Atoi(pidstr)
	if err != nil {
		log.Println("文章id截取失败")
		return
	}

	PostRes, err := service.GetDetailPost(pid)
	if err != nil {
		c.HTML(http.StatusOK, "detail.html", "文章查询失败")
		return
	}
	//每一次查看文章都应该将查看次数加1
	dao.UpdateViewCount(pid, PostRes.Article.ViewCount+1)
	c.HTML(http.StatusOK, "detail.html", PostRes)

}
