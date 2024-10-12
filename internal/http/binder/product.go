package binder

import (
	"mime/multipart"
)

type ProductCreateRequest struct {
	ProductName     string                `form:"productname" validate:"required"`
	ProductCategory string                `form:"productcategory" validate:"required"`
	SellPrice       string                `form:"sellprice" validate:"required"`
	MakePrice       string                `form:"makeprice" validate:"required"`
	Pajak           string                `form:"pajak" validate:"required"`
	Image           *multipart.FileHeader `form:"image" validate:"required"`
	Description     string                `form:"description" validate:"required"`
}

type ProductUpdateRequest struct {
	ProdukId		string                `param:"id_product" validate:"required"`
	ProductName     string                `form:"productname" validate:"required"`
	ProductCategory string                `form:"productcategory" validate:"required"`
	SellPrice       string                `form:"sellprice" validate:"required"`
	MakePrice       string                `form:"makeprice" validate:"required"`
	Pajak           string                `form:"pajak" validate:"required"`
	Image           *multipart.FileHeader `form:"image" validate:"required"`
	Description     string                `form:"description" validate:"required"`
}

type ProdukDeleteRequest struct {
	ProdukId string `param:"id_product" validate:"required"`
}
