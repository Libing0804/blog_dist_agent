package views

import (
	"blog_gin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (HTMLApi) Writing(c *gin.Context) {
	//查看一下是否有id出现在url中  如果有是修改请求
	//PIDStr := c.Query("id")
	//if PIDStr != "" {
	//	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
	//	PostAllinfo := models.Post{}
	//	c.ShouldBind(&PostAllinfo)
	//	fmt.Printf("这是全部信息%#v\n", PostAllinfo)
	//}
	wr := service.Writing()
	c.HTML(http.StatusOK, "writing.html", wr)

}
