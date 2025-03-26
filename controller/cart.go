package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartController(router *gin.Engine) {
	router.GET("/carts", GetAllPrCart)
	router.POST("/cart/add", AddToCart) // เพิ่มสินค้าลงรถเข็น
}

// เส้น API ดึงข้อมูล users
func GetAllPrCart(c *gin.Context) {
	var carts []model.Cart
	result := dbconnect.DB.Find(&carts) // ดึงข้อมูลทั้งหมด

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, carts) // ส่ง JSON response
}

// เพิ่มสินค้าลงรถเข็น
func AddToCart(c *gin.Context) {
	var req struct {
		CustomerID int    `json:"customer_id"`
		CartName   string `json:"cart_name"`
		ProductID  int    `json:"product_id"`
		Quantity   int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cart model.Cart
	if err := dbconnect.DB.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).First(&cart).Error; err != nil {
		// ถ้าไม่มีรถเข็น ให้สร้างใหม่
		cart = model.Cart{CustomerID: req.CustomerID, CartName: req.CartName}
		dbconnect.DB.Create(&cart)
	}

	var cartItem model.CartItem
	if err := dbconnect.DB.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error; err == nil {
		// ถ้ามีสินค้าอยู่แล้ว ให้เพิ่มจำนวน
		cartItem.Quantity += req.Quantity
		dbconnect.DB.Save(&cartItem)
	} else {
		// ถ้าไม่มี ให้เพิ่มรายการใหม่
		cartItem = model.CartItem{CartID: cart.CartID, ProductID: req.ProductID, Quantity: req.Quantity}
		dbconnect.DB.Create(&cartItem)
	}

	c.JSON(http.StatusOK, gin.H{"message": "เพิ่มสินค้าลงรถเข็นสำเร็จ"})
}
