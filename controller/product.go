package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ProductController(router *gin.Engine) {
	router.GET("/products", GetAllProduct)
	router.GET("/products/search", SearchProduct) //
}

// เส้น API ดึงข้อมูล users
func GetAllProduct(c *gin.Context) {
	var products []model.Product
	result := dbconnect.DB.Find(&products) // ดึงข้อมูลทั้งหมด

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products) // ส่ง JSON response
}

// ฟังก์ชันค้นหาสินค้า ตามคำค้น (description) และช่วงราคา
func SearchProduct(c *gin.Context) {
	query := c.Query("q")               // รับคำค้นหา (description)
	minPriceStr := c.Query("min_price") // รับช่วงราคาต่ำสุด
	maxPriceStr := c.Query("max_price") // รับช่วงราคาสูงสุด

	var products []model.Product
	db := dbconnect.DB

	// ค้นหาตาม description ถ้ามีการส่งค่า q มา
	if query != "" {
		db = db.Where("description LIKE ?", "%"+query+"%")
	}

	//ค้นหาตามช่วงราคา (min_price, max_price) ถ้ามีการส่งค่าเข้ามา
	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err == nil {
			db = db.Where("price >= ?", minPrice)
		}
	}

	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err == nil {
			db = db.Where("price <= ?", maxPrice)
		}
	}

	// ดึงข้อมูลสินค้าที่ตรงกับเงื่อนไข
	result := db.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
