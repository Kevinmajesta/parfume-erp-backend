package binder

type MoCreateRequest struct {
	ProductId    string `json:"id_product" validate:"required"`
	BomId        string `json:"id_bom" validate:"required"`
	Qtytoproduce string `json:"qtytoproduce" validate:"required"`
}

type UpdateMoStatusRequest struct {
	MoId string `json:"id_mo" validate:"required"`
}

type MoDeleteRequest struct {
	MoId string `param:"id_mo" validate:"required"`
}
