package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartController(router *gin.Engine) {
	router.GET("/carts", GetAllPrCart)
}

// เส้น API ดึงข้อมูล users
func GetAllPrCart(c *gin.Context) {
	var carts []model.Product
	result := dbconnect.DB.Find(&carts) // ดึงข้อมูลทั้งหมด

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, carts) // ส่ง JSON response
}
