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
	err := db.First(&p, p.ID).Error
	if err != nil {
		log.Println("Get product by id failed")
		return nil, err
	}
	return &p, nil
}

func (p Product) SelectProductIdList(db *gorm.DB, pageIndex, pageSize int) ([]uint, error) {
	var ids []uint
	err := db.Table(constant.PRODUCT_TAB).Select("id").Offset(pageIndex).Limit(pageSize).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
	//:= db.Raw("select id from "+constant.PRODUCT_TAB+" OFFSET ? LIMIT ?", pageIndex, pageSize).Scan(ids).
}

func (p Product) TableName() string {
	return constant.PRODUCT_TAB
}
