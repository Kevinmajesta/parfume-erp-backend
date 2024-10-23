package entity

import (
	"fmt"
)

type Materials struct {
	MaterialId       string `json:"id_material" gorm:"column:id_material;primaryKey"`
	Materialname     string `form:"materialname"`
	Materialcategory string `form:"materialcategory"`
	Sellprice        string `form:"sellprice"`
	Makeprice        string `form:"makeprice"`
	Unit             string `form:"unit"`
	Description      string `form:"description"`
	Image            string `form:"image"`
	Auditable
}

func generateMaterialsId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "MTR-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("MTR-%05d", newNumber)
}

func NewMaterials(lastId, materialname, materialcategory,
	sellprice, makeprice, unit, description, image string) *Materials {
	return &Materials{
		MaterialId:       generateMaterialsId(lastId),
		Materialname:     materialname,
		Materialcategory: materialcategory,
		Sellprice:        sellprice,
		Makeprice:        makeprice,
		Unit:             unit,
		Description:      description,
		Image:            image,
		Auditable:        NewAuditable(),
	}
}

func UpdateMaterials(id_material, materialname, materialcategory,
	sellprice, makeprice, unit, description string) *Materials {
	return &Materials{
		MaterialId:       id_material,
		Materialname:     materialname,
		Materialcategory: materialcategory,
		Sellprice:        sellprice,
		Makeprice:        makeprice,
		Unit:             unit,
		Description:      description,
		Auditable:        UpdateAuditable(),
	}
}
