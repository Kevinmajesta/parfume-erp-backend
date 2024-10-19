package binder


type BOMCreateRequest struct {
	IdBom            string            `json:"id_bom"` 
	IdProduct        string            `json:"id_product"`
	ProductName      string            `json:"productname"`
	ProductReference string            `json:"productpreference"`
	Quantity         string            `json:"quantity"`
	Materials        []MaterialRequest `json:"materials"`
}

type BOMUpdateRequest struct {
    IdBom            string            `json:"id_bom"`          // ID of the BoM to update
    IdProduct        string            `json:"id_product"`      // ID of the product in the BoM
    ProductName      string            `json:"productname"`     // Name of the product
    ProductReference string            `json:"productpreference"` // Reference for the product
    Quantity         string            `json:"quantity"`        // Quantity of the product
    Materials        []MaterialRequest `json:"materials"`       // List of materials for the BoM
}

type MaterialRequest struct {
	IdBomMaterial string `json:"id_bommaterial"`
	IdMaterial   string `json:"id_material"`
	MaterialName string `json:"material_name"`
	Quantity     string `json:"quantity"`
	Unit         string `json:"unit"`
}



type BomDeleteRequest struct {
	BomId string `param:"id_bom" validate:"required"`
}

