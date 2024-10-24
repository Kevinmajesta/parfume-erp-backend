package entity

import "fmt"

type Mos struct {
	MoId         string `json:"id_mo" gorm:"column:id_mo;primaryKey"`
	ProductId    string `json:"id_product" gorm:"column:id_product"`
	BomId        string `json:"id_bom" gorm:"column:id_bom"`
	Qtytoproduce string `json:"qtytoproduce"`
	Status       string `json:"status"`
	Auditable
}

func generateMosId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "MO-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("MO-%05d", newNumber)
}

func NewMos(lastId, id_product, id_bom, qtytoproduce string) *Mos {
	return &Mos{
		MoId:         generateMosId(lastId), 
		ProductId:    id_product,
		BomId:        id_bom,
		Qtytoproduce: qtytoproduce,
		Status:       "draft", 
		Auditable:    NewAuditable(),
	}
}
