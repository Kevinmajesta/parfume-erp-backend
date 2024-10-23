package entity

import (
	"fmt"
)

type Products struct {
	ProdukId        string `json:"id_product" gorm:"column:id_product;primaryKey"`
	Productname     string `form:"productname"`
	Productcategory string `form:"productcategory"`
	Sellprice       string `form:"sellprice"`
	Makeprice       string `form:"makeprice"`
	Pajak           string `form:"pajak"`
	Description     string `form:"description"`
	Variant         string `form:"variant" gorm:"column:variant"` 
	Image           string `form:"image"`
	Auditable
}

func generateProductId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "PRF-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("PRF-%05d", newNumber)
}

func NewProduct(lastId, productname, productcategory,
	sellprice, makeprice, pajak, description, image, variant string) *Products {
	return &Products{
		ProdukId:        generateProductId(lastId),
		Productname:     productname,
		Productcategory: productcategory,
		Sellprice:       sellprice,
		Makeprice:       makeprice,
		Pajak:           pajak,
		Description:     description,
		Image:           image,
		Variant:         variant,
		Auditable:       NewAuditable(),
	}
}

func UpdateProduct(id_product, productname, productcategory,
	sellprice, makeprice, pajak, description, variant  string) *Products {
	return &Products{
		ProdukId:        id_product,
		Productname:     productname,
		Productcategory: productcategory,
		Sellprice:       sellprice,
		Makeprice:       makeprice,
		Pajak:           pajak,
		Description:     description,
		Variant:         variant,
		Auditable:       UpdateAuditable(),
	}
}
