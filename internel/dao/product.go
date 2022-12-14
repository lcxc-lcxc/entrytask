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

// GetProductBriefList 获取产品列表所调用的方法
func (d *Dao) GetProductBriefList(PageIndex int, PageSize int) ([]ProductBrief, error) {
	var productBriefList []ProductBrief

	p := model.Product{}
	productList, err := p.SelectProductList(d.engine, PageIndex, PageSize)
	//
	//startTime := time.Now()
	//productIds, err := p.SelectProductIdList(d.engine, PageIndex, PageSize)
	//if err != nil {
	//	return nil, err
	//}
	//dur := time.Since(startTime)
	//fmt.Println(dur)
	//
	//productList, err := d.GetProductBriefListCache(productIds...)

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

// GetProductCount 获取产品总数
func (d *Dao) GetProductCount() (int, error) {
	p := model.Product{}
	return p.SelectProductCount(d.engine)
}

// GetProductSearch 查询商品
func (d *Dao) GetProductSearch(searchBy string) ([]ProductBrief, error) {
	var productBriefList []ProductBrief

	productSearch, err := d.GetProductSearchCache(searchBy)
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

// GetProductSearchCache 封装了查询商品时的缓存模块
func (d *Dao) GetProductSearchCache(searchBy string) ([]model.Product, error) {
	loadFunction := func(ctx context.Context, key any) ([]model.Product, error) {
		log.Println("get product search cache failed , getting data from database ")
		searchByRedisKey, _ := key.(string)
		searchBy := utils.ConvertRedisKeyToString(searchByRedisKey)

		p := model.Product{}
		productSearch, err := p.SelectProductSearch(d.engine, searchBy)

		if err != nil {
			return nil, err
		}
		return productSearch, nil
	}
	loadableCache := cache.NewLoadableCache[[]model.Product](loadFunction, d.RedisClient, time.Hour)
	productSearch, err := loadableCache.Get(context.Background(), utils.ConvertStringToRedisKey(constant.SEARCH_BY, searchBy))
	if err != nil {
		return nil, err
	}
	return productSearch, nil

}

// GetProductDetail 获取产品的详细信息：title, descriptions, product categories, product photos, and list of comments.
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

// GetProductCache 封装了获取产品表的缓存模块
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
		return product, nil
	}

	loadableCache := cache.NewLoadableCache[*model.Product](loadFunction, d.RedisClient, time.Hour)
	product, err := loadableCache.Get(context.Background(), utils.ConvertUintIdToRedisKey(constant.PRODUCT_ID, productId))

	if err != nil {
		return nil, err
	}
	return product, nil

}

// GetProductCommentInfoListCache 封装了获取产品评论的缓存模块
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

// GetProductBriefListCache
// 该函数是对产品列表对缓存模块封装：
// 逻辑是这样的 ：
// 	1. 先用pageIndex和pageSize从数据库获取到对应的id list
//  2. 后调用此方法
//  3. 使用redis的mget获取到缓存的内容，如果不存在的键，那就使用loadFunction从数据库里面获取。
//  测试了一下，使用此方法，并没有带来性能的显著提升，故没使用。
func (d *Dao) GetProductBriefListCache(productIds ...uint) ([]*model.Product, error) {

	// 从数据库获取model.Product的函数
	loadFunction := func(ctx context.Context, key any) (*model.Product, error) {
		log.Println("get product comment cache failed , getting data from database ")
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
		return product, nil
	}
	loadableCache := cache.NewLoadableCache[*model.Product](loadFunction, d.RedisClient, time.Hour)
	var keys []string
	for _, productId := range productIds {
		keys = append(keys, utils.ConvertUintIdToRedisKey(constant.PRODUCT_ID, productId))
	}
	products, err := loadableCache.MGet(context.Background(), keys...)
	if err != nil {
		return nil, err
	}
	return products, nil

}
