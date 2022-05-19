package main

import (
	"blog_gin/router"
)

func main() {
	ru := router.Router()
	ru.Run("0.0.0.0:8000")
}
