package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomerController(router *gin.Engine) {
	router.GET("/customers", GetAllCustomer)
	router.POST("/customer/login", LoginCustomer)
	router.POST("/customer/change-password", ChangePassword)
}

// เส้น API ดึงข้อมูล users
func GetAllCustomer(c *gin.Context) {
	var customers []model.Customer
	result := dbconnect.DB.Find(&customers) // ดึงข้อมูลทั้งหมด

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, customers) // ส่ง JSON response
}

// Login API
func LoginCustomer(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// ตรวจสอบว่าได้รับข้อมูลครบหรือไม่
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบ"})
		return
	}

	var customer model.Customer
	result := dbconnect.DB.Where("email = ?", loginData.Email).First(&customer)

	// ตรวจสอบว่าอีเมลถูกต้องหรือไม่
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "อีเมลหรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	// ตรวจสอบรหัสผ่าน
	if customer.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสผ่านไม่ถูกต้อง"})
		return
	}

	// สร้างโครงสร้างข้อมูลใหม่สำหรับ Response (ไม่รวม Password)
	responseData := struct {
		CustomerID  uint   `json:"CustomerID"`
		FirstName   string `json:"FirstName"`
		LastName    string `json:"LastName"`
		Email       string `json:"Email"`
		PhoneNumber string `json:"PhoneNumber"`
		Address     string `json:"Address"`
		CreatedAt   string `json:"CreatedAt"`
	}{
		CustomerID:  uint(customer.CustomerID),
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Address:     customer.Address,
		CreatedAt:   customer.CreatedAt.String(),
	}

	// ส่ง Response กลับโดยไม่รวม Password
	c.JSON(http.StatusOK, gin.H{"message": "เข้าสู่ระบบสำเร็จ", "customer": responseData})
}

// API สำหรับเปลี่ยนรหัสผ่าน
func ChangePassword(c *gin.Context) {
	// รับข้อมูลจาก Request JSON
	var data struct {
		Email       string `json:"email" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// ตรวจสอบข้อมูลที่ส่งมาว่าครบหรือไม่
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบ"})
		return
	}

	// ค้นหาลูกค้าจาก email
	var customer model.Customer
	result := dbconnect.DB.Where("email = ?", data.Email).First(&customer)

	// ตรวจสอบว่าพบลูกค้าหรือไม่
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลลูกค้า"})
		return
	}

	// ตรวจสอบรหัสผ่านเก่าว่าตรงกับที่เก็บไว้หรือไม่ (ไม่มีการแฮช)
	if customer.Password != data.OldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสผ่านเก่าผิด"})
		return
	}

	// ตรวจสอบว่ารหัสผ่านใหม่ต้องไม่เหมือนรหัสผ่านเก่า
	if data.OldPassword == data.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสผ่านใหม่ต้องไม่เหมือนรหัสผ่านเก่า"})
		return
	}

	// ตรวจสอบความยาวรหัสผ่านใหม่ (อย่างน้อย 6 ตัว)
	if len(data.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสผ่านใหม่ต้องมีอย่างน้อย 6 ตัวอักษร"})
		return
	}

	// อัพเดตรหัสผ่านใหม่ในฐานข้อมูล (ไม่แฮช)
	customer.Password = data.NewPassword
	if err := dbconnect.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัพเดตรหัสผ่านได้"})
		return
	}

	// ส่งข้อความตอบกลับเมื่อเปลี่ยนรหัสผ่านสำเร็จ
	c.JSON(http.StatusOK, gin.H{"message": "เปลี่ยนรหัสผ่านสำเร็จ"})
}
