package binder


type BOMCreateRequest struct {
	IdProduct        string            `json:"id_product"`
	ProductName      string            `json:"productname"`
	ProductReference string            `json:"productpreference"`
	Quantity         string            `json:"quantity"`
	Unit             string            `json:"unit"`
	Materials        []MaterialRequest `json:"materials"`
}

type MaterialRequest struct {
	IdBomMaterial string `json:"id_bommaterial"`
	IdMaterial   string `json:"id_material"`
	MaterialName string `json:"material_name"`
	Quantity     string `json:"quantity"`
	Unit         string `json:"unit"`
}


