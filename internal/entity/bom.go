package entity

import (
	"fmt"
)

type Bom struct {
	BomId             string        `json:"id_bom" gorm:"column:id_bom;primaryKey"`
	ProductId         string        `json:"id_product" gorm:"column:id_product"`
	ProductName       string        `json:"productname" gorm:"column:productname"`
	ProductPreference string        `json:"productpreference" gorm:"column:productpreference"`
	Quantity          string        `json:"quantity" gorm:"column:quantity"`
	Unit              string        `json:"unit" gorm:"column:unit"`
	Materials         []BomMaterial `json:"materials" gorm:"-"`
	Auditable
}

type BomMaterial struct {
	IdBomMaterial string `json:"id_bommaterial" gorm:"column:id_bommaterial"`
	IdMaterial    string `json:"id_material" gorm:"column:id_material"`
	BomId         string `json:"id_bom" gorm:"column:id_bom"`
	MaterialName  string `json:"material_name" gorm:"column:materialname"`
	Quantity      string `json:"quantity" gorm:"column:quantity"`
	Unit          string `json:"unit" gorm:"column:unit"`
	Auditable
}

func generateBomId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "BOM-%d", &newNumber)
		newNumber++
	}
	return fmt.Sprintf("BOM-%05d", newNumber)
}

func NewBom(lastId, productId, productName, productpreference, quantity, unit string) *Bom {
	return &Bom{
		BomId:             generateBomId(lastId),
		ProductId:         productId,
		ProductName:       productName,
		ProductPreference: productpreference,
		Quantity:          quantity,
		Unit:              unit,
		Auditable:         NewAuditable(),
	}
}
