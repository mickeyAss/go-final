package controller

import (
	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	CartItemController(router)
	CartController(router)
	ProductController(router)
	CustomerController(router)
	router.Run()
}
