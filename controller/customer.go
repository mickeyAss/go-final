package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserController(router *gin.Engine) {
	router.GET("/ping", ping)
	router.GET("/users", GetAllUsers)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Pong!!",
	})
}

// ✅ เส้น API ดึงข้อมูล users
func GetAllUsers(c *gin.Context) {
	var customers []model.Customer
	result := dbconnect.DB.Find(&customers) // ดึงข้อมูลทั้งหมด

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customers) // ส่ง JSON response
}
