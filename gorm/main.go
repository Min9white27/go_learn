package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/your_db"))
	if err != nil {
		panic("failed to connect database")
	}

	db = db.Debug()
	// 迁移 schema
	// 建表
	db.AutoMigrate(&Product{})
	if err != nil {
		return
	}

	// Create 插入
	db.Create(&Product{Code: "D42", Price: 100})

	// Read 查找
	var product Product
	db.First(&product, 1)                 // 根据整型主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	//也就是这一句话会更新Price 和 Code 两个字段
	//可以看作为 SET `Price`=200,`Code`="F42"
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}
