package controller

import (
	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	CartController(router)
	ProductController(router)
	CustomerController(router)
	router.Run()
}
