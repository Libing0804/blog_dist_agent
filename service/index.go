package service

import (
	"blog_gin/config"
	"blog_gin/dao"
	"blog_gin/daors"
	"blog_gin/models"
	"encoding/json"
	"html/template"
	"log"
)

func GetAllIndexInfo(slug string, page, pageSize int) (*models.HomeResponse, error) {
	categorys, err := dao.GetAllCategory()

	if err != nil {
		log.Println(err)
		return nil, err

	}
	var posts []models.Post
	var postMores []models.PostMore
	if slug == "" {

		posts, err = dao.GetPostPage(page, pageSize)
		//将文章数据缓存进入redis
		//先清空数据
		_, errRedis := daors.Deleteredis()
		if errRedis != nil {
			log.Println("redis数据清空失败")
		}
		var postString = make([]string, 0)
		for _, post := range posts {
			byteTemp, errRedis := json.Marshal(post)
			if errRedis != nil {
				log.Println("将文章缓存进入redis数据库时，序列化失败！！！")
			}
			postString = append(postString, string(byteTemp))
		}
		//放入数据库
		daors.Add(postString)

	} else {

		posts, err = dao.GetPostPageBySlug(slug, page, pageSize)
	}

	for _, post := range posts {
		categoryName, err := dao.GetCategoryById(post.CategoryId)
		if err != nil {
			log.Println("分类名查询失败")
		}
		UserName, err := dao.GetUserNameById(post.UserId)
		if err != nil {
			log.Println("用户名查询失败")
		}
		//页面展示的时候太乱了  view展示页面只显示部分内容
		contentTemp := []rune(post.Content)
		if len(contentTemp) > 100 {

			post.Content = string(contentTemp[:100])
		}
		var postMore models.PostMore
		postMore.Pid = post.Pid
		postMore.Title = post.Title

		postMore.Slug = post.Slug
		postMore.Content = template.HTML(post.Content)
		postMore.CategoryName = categoryName
		postMore.CategoryId = post.CategoryId
		postMore.UserId = post.UserId
		postMore.UserName = UserName
		postMore.ViewCount = post.ViewCount
		postMore.Type = post.Type
		postMore.CreateAt = post.CreateAt.Format("2006-01-02 15:04:05")
		postMore.UpdateAt = post.UpdateAt.Format("2006-01-02 15:04:05")
		postMores = append(postMores, postMore)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var total int
	//查询总文章数 和页数
	if slug == "" {
		total = dao.CountGetAllPost()
	} else {
		total = dao.CountGetAllPostBySlug(slug)
	}

	allPages := (total-1)/10 + 1
	var pages = []int{}
	for i := 0; i < allPages; i++ {
		pages = append(pages, i+1)
	}
	var hr = &models.HomeResponse{
		config.Cfg.Viewer,
		categorys,
		postMores,
		total,
		page,
		pages,
		page != allPages,
	}
	return hr, nil

}

//查询主页面的所有信息
func GetAllIndexInfoRedis(page, pageSize int) (*models.HomeResponse, error) {
	categorys, err := dao.GetAllCategory()

	if err != nil {
		log.Println(err)
		return nil, err

	}
	var posts []models.Post
	var postMores []models.PostMore
	posts, err = daors.GetPostAll(page, pageSize)
	for _, post := range posts {
		categoryName, err := dao.GetCategoryById(post.CategoryId)
		if err != nil {
			log.Println("分类名查询失败")
		}
		UserName, err := dao.GetUserNameById(post.UserId)
		if err != nil {
			log.Println("用户名查询失败")
		}
		//页面展示的时候太乱了  view展示页面只显示部分内容
		contentTemp := []rune(post.Content)
		if len(contentTemp) > 100 {

			post.Content = string(contentTemp[:100])
		}
		var postMore models.PostMore
		postMore.Pid = post.Pid
		postMore.Title = post.Title

		postMore.Slug = post.Slug
		postMore.Content = template.HTML(post.Content)
		postMore.CategoryName = categoryName
		postMore.CategoryId = post.CategoryId
		postMore.UserId = post.UserId
		postMore.UserName = UserName
		postMore.ViewCount = post.ViewCount
		postMore.Type = post.Type
		postMore.CreateAt = post.CreateAt.Format("2006-01-02 15:04:05")
		postMore.UpdateAt = post.UpdateAt.Format("2006-01-02 15:04:05")
		postMores = append(postMores, postMore)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var total int
	//查询总文章数 和页数

	total = int(daors.DBrs.LLen("postsall").Val())

	allPages := (total-1)/10 + 1
	var pages = []int{}
	for i := 0; i < allPages; i++ {
		pages = append(pages, i+1)
	}
	var hr = &models.HomeResponse{
		config.Cfg.Viewer,
		categorys,
		postMores,
		total,
		page,
		pages,
		page != allPages,
	}
	return hr, nil

}
