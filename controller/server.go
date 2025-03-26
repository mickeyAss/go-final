package controller

import (
	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	ProductController(router)
	CustomerController(router)
	router.Run()
}
