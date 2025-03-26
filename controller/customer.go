package controller

import (
	"go-basic/dbconnect"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	var data struct {
		Email       string `json:"email" binding:"required"`
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// ตรวจสอบว่าได้รับข้อมูลครบหรือไม่
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบ"})
		return
	}

	// ค้นหาลูกค้าจาก email
	var customer model.Customer
	result := dbconnect.DB.Where("email = ?", data.Email).First(&customer)

	// ตรวจสอบว่าเจอลูกค้าหรือไม่
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลลูกค้า"})
		return
	}

	// ตรวจสอบรหัสผ่านเก่ากับที่เก็บไว้ในฐานข้อมูล (bcrypt)
	err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(data.OldPassword))
	if err != nil {
		// หากมีข้อผิดพลาดในการเปรียบเทียบรหัสผ่าน
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสผ่านเก่าผิด"})
		return
	}

	// เข้ารหัสรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเข้ารหัสรหัสผ่านใหม่ได้"})
		return
	}

	// อัพเดตรหัสผ่านใหม่ในฐานข้อมูล
	customer.Password = string(hashedPassword)
	if err := dbconnect.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัพเดตข้อมูลได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "เปลี่ยนรหัสผ่านสำเร็จ"})
}
