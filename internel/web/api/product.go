package api

import (
	"entrytask/internel/constant"
	http_service "entrytask/internel/service/http-service"
	"entrytask/internel/web/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Product struct {
}

func NewProduct() *Product {
	return &Product{}
}

func (p *Product) List(c *gin.Context) {
	resp := response.NewResponse(c)
	var param http_service.ProductListRequest

	param.PageIndex, param.PageSize = cast.ToInt(c.Query(constant.PAGE_INDEX)), cast.ToInt(c.Query(constant.PAGE_SIZE))
	if param.PageIndex <= 0 {
		param.PageIndex = 1
	}
	if param.PageSize <= 0 || param.PageSize > 20 {
		param.PageSize = 10
	}

	service := http_service.NewService(c.Request.Context())
	productListResponse, err := service.ProductList(&param)
	if err != nil {
		resp.ResponseError(constant.ProductListGetFailed.GetRetCode(), err.Error())
		return
	}
	resp.ResponseOK(productListResponse)
}

func (p *Product) Search(c *gin.Context) {
	resp := response.NewResponse(c)

	searchBy := c.Query(constant.SEARCH_BY)
	if searchBy == "" {
		resp.ResponseError(constant.InvalidParams.GetRetCode(), "校验失败（请输入需要搜索的字段）")
		return
	}
	service := http_service.NewService(c.Request.Context())
	search, err := service.ProductSearch(&http_service.ProductSearchRequest{SearchBy: searchBy})
	if err != nil {
		resp.ResponseError(constant.ProductSearchFailed.GetRetCode(), err.Error())
		return
	}
	resp.ResponseOK(search)

}

func (p *Product) Detail(c *gin.Context) {
	resp := response.NewResponse(c)
	productId := cast.ToUint(c.Param(constant.PRODUCT_ID))
	if productId == 0 {
		resp.ResponseError(constant.InvalidParams.GetRetCode(), "校验失败")
		return
	}

	service := http_service.NewService(c.Request.Context())
	productDetail, err := service.ProductDetail(&http_service.ProductDetailRequest{ProductId: productId})
	if err != nil {
		resp.ResponseError(constant.ProductDetailGetFailed.GetRetCode(), err.Error())
		return
	}

	resp.ResponseOK(productDetail)
}
