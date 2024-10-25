package binder

type MoCreateRequest struct {
	ProductId    string `json:"id_product" validate:"required"`
	BomId        string `json:"id_bom" validate:"required"`
	Qtytoproduce string `json:"qtytoproduce" validate:"required"`
}

