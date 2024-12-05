package binder

import (
	"mime/multipart"
)

type MaterialCreateRequest struct {
	MaterialName     string                `form:"materialname" validate:"required"`
	MaterialCategory string                `form:"materialcategory" validate:"required"`
	SellPrice        string                `form:"sellprice" validate:"required"`
	MakePrice        string                `form:"makeprice" validate:"required"`
	Unit             string                `form:"unit" validate:"required"`
	Image            *multipart.FileHeader `form:"image" validate:"required"`
	Description      string                `form:"description" validate:"required"`
}

type MaterialUpdateRequest struct {
	MaterialId        string                `param:"id_material" validate:"required"`
	MaterialtName     string                `form:"materialname" validate:"required"`
	MaterialtCategory string                `form:"materialcategory" validate:"required"`
	SellPrice         string                `form:"sellprice" validate:"required"`
	MakePrice         string                `form:"makeprice" validate:"required"`
	Unit              string                `form:"unit" validate:"required"`
	Image             *multipart.FileHeader `form:"image" validate:"required"`
	Description       string                `form:"description" validate:"required"`
}

type MaterialDeleteRequest struct {
	MaterialId string `param:"id_material" validate:"required"`
}

type MinQtyMaterial struct {
	MaterialName string  `json:"MaterialName" validate:"required"` // Update to match incoming JSON
	Qty          float64 `json:"Qty" validate:"required,numeric"`  // Update to match incoming JSON
}
