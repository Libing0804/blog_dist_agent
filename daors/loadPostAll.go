package daors

import (
	"blog_gin/models"
	"encoding/json"
	"log"
	"time"
)

//缓存文章
func Add(postString []string) {
	// 从右边放入元素
	DBrs.RPush("postsall", postString)
	//缓存所有文章数据并且过期时间为2小时
	DBrs.Expire("postsall", 7200*time.Second)
}

//每次增删改查都应该删除redis中关于主页全文的缓存
func Deleteredis() (n int64, err error) {

	n, err = DBrs.Del("postsall").Result()
	return
}

//查询文章信息
func GetPostAll(page, pageSize int) ([]models.Post, error) {
	page = (page - 1) * pageSize

	postsRes := DBrs.LRange("postsall", int64(page), int64(page+pageSize)).Val()
	var ret = make([]models.Post, 0)
	for _, v := range postsRes {
		t := models.Post{}
		v1 := []byte(v)
		err := json.Unmarshal(v1, &t)
		if err != nil {
			log.Println("redis中的文章反序列化失败")
			return nil, err
		}
		ret = append(ret, t)
	}
	return ret, nil
}
