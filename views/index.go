package views

import (
	"blog_gin/daors"
	"blog_gin/models"
	"blog_gin/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (*HTMLApi) Index(c *gin.Context) {

	//获取表单
	pageStr := c.Param("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	//每页显示的数量
	pageSize := 10
	//去数据库获取需要的数据
	pathurl := c.Request.RequestURI
	slug := strings.TrimPrefix(pathurl, "/")

	var err error
	var hr *models.HomeResponse
	//先查询redis
	n, _ := daors.DBrs.Exists("postsall").Result()
	if n > 0 {
		//存在
		fmt.Println("*************走的redis缓存**********************")
		tim := daors.DBrs.TTL("postsall").Val()
		fmt.Printf("缓存过期时间剩余时间%v\n", tim)

		hr, err = service.GetAllIndexInfoRedis(page, pageSize)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusOK, "index.html", errors.New("系统错误"))
			return
		}
	} else {
		//redis不存在才会查询mysql并且将数据缓存进入redis
		hr, err = service.GetAllIndexInfo(slug, page, pageSize)

		if err != nil {
			log.Println(err)
			c.HTML(http.StatusOK, "index.html", errors.New("系统错误"))
			return
		}

	}
	c.HTML(http.StatusOK, "index.html", hr)

}
