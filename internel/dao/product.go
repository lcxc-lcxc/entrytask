package dao

import (
	"entrytask/internel/model"
	"gorm.io/gorm"
)

type ProductBrief struct {
	ProductId         uint   `json:"product_id"`
	ProductName       string `json:"product_name"`
	CategoryName      string `json:"category_name"`
	ProductPictureUrl string `json:"product_picture_url"`
}

type ProductDetail struct {
	ProductId          uint           `json:"product_id"`
	ProductName        string         `json:"product_name"`
	ProductDescription string         `json:"product_description"`
	CommentList        []CommentBrief `json:"comment_list,omitempty"`
}

func (d *Dao) GetProductBriefList(PageIndex int, PageSize int) ([]ProductBrief, error) {
	var productBriefList []ProductBrief

	p := model.Product{}
	productList, err := p.SelectProductList(d.engine, PageIndex, PageSize)
	if err != nil {
		return nil, err
	}
	for _, product := range productList {
		productBriefList = append(productBriefList, ProductBrief{
			ProductId:         product.Model.ID,
			ProductName:       product.ProductName,
			CategoryName:      product.CategoryName,
			ProductPictureUrl: product.PictureUrl,
		})
	}
	return productBriefList, nil

}

func (d *Dao) GetProductCount() (int, error) {
	p := model.Product{}
	return p.SelectProductCount(d.engine)
}

func (d *Dao) GetProductSearch(searchBy string) ([]ProductBrief, error) {
	var productBriefList []ProductBrief
	p := model.Product{}
	productSearch, err := p.SelectProductSearch(d.engine, searchBy)
	if err != nil {
		return nil, err
	}
	for _, product := range productSearch {
		productBriefList = append(productBriefList, ProductBrief{
			ProductId:         product.Model.ID,
			ProductName:       product.ProductName,
			CategoryName:      product.CategoryName,
			ProductPictureUrl: product.PictureUrl,
		})
	}
	return productBriefList, nil

}

func (d *Dao) GetProductDetail(productId uint) (*ProductDetail, error) {

	// 1 获取product表内的所有信息
	p := model.Product{
		Model: gorm.Model{
			ID: productId,
		},
	}
	product, err := p.SelectProduct(d.engine)
	if err != nil {
		return nil, err
	}
	// 2 获取product的评论信息
	c := model.CommentInfo{}

	commentInfoList, err := c.SelectCommentInfoList(d.engine, productId)
	if err != nil {
		return nil, err
	}

	// 3 封装数据
	productDetail := &ProductDetail{
		ProductId:          product.ID,
		ProductName:        product.ProductName,
		ProductDescription: product.Description,
	}

	var productCommentList []CommentBrief
	for _, commentInfo := range commentInfoList {
		productCommentList = append(productCommentList, CommentBrief{
			CommentId: commentInfo.ID,
			FromName:  commentInfo.FromName,
			Content:   commentInfo.Content,
			CreatedAt: commentInfo.CreatedAt,
		})
	}
	productDetail.CommentList = productCommentList

	return productDetail, nil

}
