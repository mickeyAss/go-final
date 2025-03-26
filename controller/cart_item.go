package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartItemController(router *gin.Engine) {
	router.GET("/cartItems", GetAllPrCart)
}

// เส้น API ดึงข้อมูล users
func GetAllPrCartItem(c *gin.Context) {
	var cartItems []model.CartItem
	result := dbconnect.DB.Find(&cartItems) // ดึงข้อมูลทั้งหมด

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, cartItems) // ส่ง JSON response
}
