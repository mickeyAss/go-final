package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CartItemController(router *gin.Engine) {
	router.GET("/cartItems", GetCustomerCarts)
}

// ดึงข้อมูลรถเข็นทั้งหมดของลูกค้าพร้อมรายละเอียดสินค้า
func GetCustomerCarts(c *gin.Context) {
	customerID := c.Query("customer_id")
	var carts []model.Cart

	// ดึงข้อมูลรถเข็นของลูกค้า
	if err := dbconnect.DB.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cartDetails []gin.H
	for _, cart := range carts {
		var cartItems []model.CartItem
		dbconnect.DB.Where("cart_id = ?", cart.CartID).Find(&cartItems)

		var items []gin.H

		for _, item := range cartItems {
			var product model.Product
			dbconnect.DB.Where("product_id = ?", item.ProductID).First(&product)

			items = append(items, gin.H{
				"product_id":   product.ProductID,
				"product_name": product.ProductName,
				"description":  product.Description,
				"price":        product.Price,
				"quantity":     item.Quantity,
			})
		}

		cartDetails = append(cartDetails, gin.H{
			"cart_id":   cart.CartID,
			"cart_name": cart.CartName,
			"items":     items,
		})
	}

	c.JSON(http.StatusOK, cartDetails)
}
