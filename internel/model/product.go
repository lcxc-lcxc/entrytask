package model

import (
	"entrytask/internel/constant"
	"gorm.io/gorm"
	"log"
)

//与表中字段完全映射
type Product struct {
	ProductName  string
	CategoryName string
	Description  string
	PictureUrl   string
	gorm.Model
}

func (p Product) SelectProductList(db *gorm.DB, PageIndex int, PageSize int) ([]Product, error) {
	var productList []Product
	err := db.Model(&p).Offset((PageIndex - 1) * PageSize).Limit(PageSize).Find(&productList).Error
	if err != nil {
		log.Println("Get Product List Failed")
		return nil, err
	}
	return productList, nil
}

func (p Product) SelectProductCount(db *gorm.DB) (int, error) {
	var count int64
	err := db.Model(&p).Count(&count).Error
	if err != nil {
		log.Println("Get Product count Failed")
		return 0, err
	}
	return int(count), nil
}

func (p Product) SelectProductSearch(db *gorm.DB, searchBy string) ([]Product, error) {
	var productList []Product
	err := db.Model(&p).Where("product_name = ?", searchBy).Or("category_name = ?", searchBy).Find(&productList).Error

	if err != nil {
		log.Println("Search product by name or category failed")
		return nil, err
	}
	return productList, nil
}

func (p Product) SelectProduct(db *gorm.DB) (*Product, error) {
	var product Product
	err := db.Model(&p).First(&product).Error
	if err != nil {
		log.Println("Get product by id failed")
		return nil, err
	}
	return &product, nil
}

func (p Product) TableName() string {
	return constant.PRODUCT_TAB
}
