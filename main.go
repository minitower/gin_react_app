package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static",
		"/var/www/diary_helper/static")
	r.LoadHTMLFiles("/var/www/diary_helper/index.html")
	BuildRouter(r)
	gin.Logger()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Run server error: " + err.Error() + "!")
	}
}
