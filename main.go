package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello, AWS App Runner!")
	})

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("サーバー起動に失敗", err)
	}
}
