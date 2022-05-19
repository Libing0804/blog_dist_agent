package api

import (
	"blog_gin/dao"
	"blog_gin/daors"
	"blog_gin/models"
	"blog_gin/service"
	"blog_gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SaveP struct {
	CategoryId string `json:"categoryId" uri:"categoryId" form:"categoryId" param:"categoryId"`
	Content    string `json:"content" uri:"content" form:"content" param:"content"`
	Markdown   string `json:"markdown" uri:"markdown" form:"markdown" param:"markdown"`
	Slug       string `json:"slug" uri:"slug" form:"slug" param:"slug"`
	Title      string `json:"title" uri:"title" form:"title" param:"title"`
	Type       string `json:"type" uri:"type" form:"type" param:"type"`
}

func (*Api) DeletePost(c *gin.Context) {
	//redis这个时候应该清空，因为缓存的数据已经变化了
	_, errRedis := daors.Deleteredis()
	if errRedis != nil {
		log.Println("清空redis失败")
	}

	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Println("获取文章pid失败")
		return
	}
	//删除之前应该验证是不是自己的文章  是自己的才能删除
	//获取用户id判断是否登录
	token := c.Request.Header.Get("Authorization")

	_, claim, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"error": "",
			"data":  "你没有登录！！！请先进行登录。",
		})
		return
	}

	userId := claim.Uid
	//查询要删除的id是不是自己的文章，否则删不掉 如果是管理员也可以删除
	p, _ := dao.GetPostById(pid)

	if p.UserId == userId || userId == 1 {
		dao.DeletePost(pid)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":  403,
			"error": "这不是你的文章，没用权限删除，请查看清楚！！！",
			"data":  "",
		})
		return
	}

	post, err := dao.GetAllPost()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  -999,
			"error": err,
			"data":  "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"error": "",
		"data":  post,
	})
}
func (*Api) GetPost(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Println("获取文章pid失败")
		return
	}

	Post, err := dao.GetPostById(pid)
	if err != nil {
		c.HTML(http.StatusOK, "detail.html", "文章查询失败")
		return
	}

	c.JSON(200, gin.H{
		"code":  200,
		"error": "",
		"data":  Post,
	})

	//c.HTML(http.StatusOK, "detail.html", PostRes)
}

func (*Api) SavePost(c *gin.Context) {
	//redis这个时候应该清空，因为缓存的数据已经变化了
	_, errRedis := daors.Deleteredis()
	if errRedis != nil {
		log.Println("清空redis失败")
	}

	//获取用户id判断是否登录
	token := c.Request.Header.Get("Authorization")
	_, claim, err := utils.ParseToken(token)
	if err != nil {
		c.HTML(200, "writing.html", err)
		return
	}
	uid := claim.Uid
	//有几个键不填需要设置默认值
	var pv = SaveP{
		CategoryId: "2",
		Slug:       "default",
		Type:       "0",
	}

	c.ShouldBind(&pv)

	//	post save

	cid, _ := strconv.Atoi(pv.CategoryId)

	typeId, _ := strconv.Atoi(pv.Type)
	post := &models.Post{
		Pid:        -1,
		Title:      pv.Title,
		Slug:       pv.Slug,
		Content:    pv.Content,
		Markdown:   pv.Markdown,
		CategoryId: cid,
		UserId:     uid,
		ViewCount:  0,
		Type:       typeId,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	}

	service.SavePost(post)

	c.JSON(200, gin.H{
		"code":  200,
		"error": "这里没",
		"data":  post,
	})
}

//创建一个简单的结构体 用于获取前端的文章id
type GetPid struct {
	PostID     int    `json:"pid" form:"pid" url:"pid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Slug       string `json:"slug"`
	Type       int    `json:"type"`
	CategoryID string `json:"category_id"`
	Markdown   string `json:"markdown"`
}

func (*Api) UpdatePost(c *gin.Context) {
	//redis这个时候应该清空，因为缓存的数据已经变化了
	_, errRedis := daors.Deleteredis()
	if errRedis != nil {
		log.Println("清空redis失败")
	}

	//	post 更新
	//获取用户id判断是否登录
	token := c.Request.Header.Get("Authorization")
	_, claim, err := utils.ParseToken(token)
	if err != nil {
		c.HTML(200, "writing.html", err)
		return
	}
	//获取文章的id
	var PostBind = GetPid{}
	err = c.ShouldBind(&PostBind)
	//没写的就是默认
	if PostBind.CategoryID != "1" || PostBind.CategoryID != "2" {
		PostBind.CategoryID = "2"
	}
	if PostBind.Slug == "" {
		PostBind.Slug = "default"
	}
	if err != nil {
		fmt.Printf("更新文章时前端数据绑定失败err:%s\n", err)
	}
	postinfo, _ := dao.GetPostById(PostBind.PostID)
	uid := claim.Uid

	//验证登录的和修改的是不是同一个用户的

	if uid != postinfo.UserId {
		c.JSON(200, gin.H{
			"code":  403,
			"error": "这不是您的文章！！！，请查看清楚。",
			"data":  "",
		})
		return
	}
	cid, _ := strconv.Atoi(PostBind.CategoryID)
	post := &models.Post{
		Pid:        PostBind.PostID,
		Title:      PostBind.Title,
		Slug:       PostBind.Slug,
		Content:    PostBind.Content,
		Markdown:   PostBind.Markdown,
		CategoryId: cid,
		UserId:     uid,
		ViewCount:  postinfo.ViewCount,
		Type:       PostBind.Type,
		CreateAt:   postinfo.CreateAt,
		UpdateAt:   time.Now(),
	}
	//fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#######")
	//fmt.Printf("获取文章的id:%d\n", post.Pid)
	//fmt.Printf("获取文章的id:%d\n", post.ViewCount)
	//
	//fmt.Printf("获取文章的title:%s\n", post.Title)
	//fmt.Printf("content:%s\n", post.Content)
	//fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	service.UpdatePost(post)
	c.JSON(200, gin.H{
		"code":  200,
		"error": "这里没错",
		"data":  post,
	})
}
