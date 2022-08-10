package dao

import (
	"context"
	"entrytask/internel/constant"
	"entrytask/internel/dao/cache"
	"entrytask/internel/model"
	"entrytask/pkg/utils"
	"gorm.io/gorm"
	"log"
	"time"
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

////此方法用于缓存ProductDetail时进行序列化
//func (p *ProductDetail) MarshalBinary() (data []byte, err error) {
//	return json.Marshal(p)
//}
//
//func (p *ProductDetail) UnmarshalBinary(data []byte) error {
//	return json.Unmarshal(data, p)
//}

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

	product, err := d.GetProductCache(productId)
	if err != nil {
		return nil, err
	}

	//2 获取product的评论信息

	commentInfoList, err := d.GetProductCommentInfoListCache(productId)

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

func (d *Dao) GetProductCache(productId uint) (*model.Product, error) {
	// 1 缓存层

	loadFunction := func(ctx context.Context, key any) (*model.Product, error) {
		log.Println("get product cache failed , getting data from database ")
		productIdStr, _ := key.(string)
		productId, _ := utils.ConvertRedisKeyToUintId(productIdStr)

		p := model.Product{
			Model: gorm.Model{
				ID: productId,
			},
		}

		product, err := p.SelectProduct(d.engine)
		if err != nil {
			return nil, err
		}
		//访问数据库后一定要以 []byte返回（通过msgpack的marshal方法），否则会出错
		return product, nil
	}

	loadableCache := cache.NewLoadableCache[*model.Product](loadFunction, d.RedisClient, time.Hour)
	product, err := loadableCache.Get(context.Background(), utils.ConvertUintIdToRedisKey(constant.PRODUCT_ID, productId))

	if err != nil {
		return nil, err
	}
	return product, nil

}

func (d *Dao) GetProductCommentInfoListCache(productId uint) ([]model.CommentInfo, error) {

	// 1 从数据库获取内容的函数
	loadFunction := func(ctx context.Context, key any) ([]model.CommentInfo, error) {
		log.Println("get product comment cache failed , getting data from database ")
		productIdStr, _ := key.(string)
		productId, _ := utils.ConvertRedisKeyToUintId(productIdStr)

		c := model.CommentInfo{}

		commentInfoList, err := c.SelectCommentInfoList(d.engine, productId) //从数据库获取内容
		if err != nil {
			return nil, err
		}
		return commentInfoList, nil
	}

	loadableCache := cache.NewLoadableCache[[]model.CommentInfo](loadFunction, d.RedisClient, time.Hour)

	commentInfoList, err := loadableCache.Get(context.Background(), utils.ConvertUintIdToRedisKey(constant.COMMENT_OF_PRODUCT, productId))
	if err != nil {
		return nil, err
	}

	return commentInfoList, nil

}

func (d *Dao) GetProductBriefListCache(PageIndex int, PageSize int) ([]ProductBrief, error) {

	return nil, nil
}
