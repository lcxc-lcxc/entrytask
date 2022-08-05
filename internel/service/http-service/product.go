package http_service

import (
	"entrytask/internel/dao"
)

type ProductListRequest struct {
	PageIndex int
	PageSize  int
}

type ProductListResponse struct {
	ProductList  []dao.ProductBrief `json:"product_list"`
	ProductCount int                `json:"product_count"`
}

type ProductSearchRequest struct {
	SearchBy string
}
type ProductSearchResponse struct {
	ProductSearchList []dao.ProductBrief `json:"product_search_list"`
}

type ProductDetailRequest struct {
	ProductId uint
}

type ProductDetailResponse struct {
	ProductDetail *dao.ProductDetail
}

func (svc *Service) ProductList(request *ProductListRequest) (*ProductListResponse, error) {
	productBriefList, err := svc.dao.GetProductBriefList(request.PageIndex, request.PageSize)
	if err != nil {
		return nil, err
	}
	count, err := svc.dao.GetProductCount()
	if err != nil {
		return nil, err
	}
	return &ProductListResponse{
		ProductList:  productBriefList,
		ProductCount: count,
	}, nil
}

func (svc *Service) ProductSearch(request *ProductSearchRequest) (*ProductSearchResponse, error) {
	productSearch, err := svc.dao.GetProductSearch(request.SearchBy)
	if err != nil {
		return nil, err
	}
	return &ProductSearchResponse{
		ProductSearchList: productSearch,
	}, nil

}

func (svc *Service) ProductDetail(request *ProductDetailRequest) (*ProductDetailResponse, error) {

	productDetail, err := svc.dao.GetProductDetail(request.ProductId)
	if err != nil {
		return nil, err
	}
	return &ProductDetailResponse{
		ProductDetail: productDetail,
	}, nil
}

func GetProductDetailDao() {

}
