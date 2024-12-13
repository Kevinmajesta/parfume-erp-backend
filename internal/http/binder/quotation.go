package binder

type QUOCreateRequest struct {
	QuotationsId string           `json:"id_quotation"`
	OrderDate    string           `json:"order_date"`
	Status       string           `json:"status"`
	CostumerId   string           `json:"id_costumer"`
	Products     []ProductRequest `json:"products"`
}

type QUOUpdateRequest struct {
	QuotationsId string              `param:"id_quotation"`
	ProductId    string              `json:"id_product"`
	OrderDate    string              `json:"order_date"`
	CostumerId   string              `json:"id_costumer"`
	Status       string              `json:"status"`
	Products     []QUOProductRequest `json:"products"`
}

type QUOProductRequest struct {
	QuotationsProductId string `json:"id_quotationsproduct"`
	ProductId           string `json:"id_product"`
	CostumerId          string `json:"id_costumer"`
	ProductName         string `json:"productname"`
	Quantity            string `json:"quantity"`
	UnitPrice           string `json:"unitprice"`
	Tax                 string `json:"tax"`
	Subtotal            string `json:"subtotal"`
}

type UpdateQuoStatusRequest struct {
	QuotationsId string `param:"id_quotation" validate:"required"`
}

type QuoDeleteRequest struct {
	QuotationsId string `param:"id_quotation" validate:"required"`
}
