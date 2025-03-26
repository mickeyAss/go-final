package main

import (
	"fmt"
	"go-basic/controller"
	"go-basic/dbconnect"
	"go-basic/model"
)

func main() {
	db := dbconnect.ConnectDB() // เรียกใช้ฟังก์ชันเชื่อมต่อฐานข้อมูล

	customer := []model.Customer{}
	result := db.Find(&customer)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(customer)

	controller.StartServer()
}
