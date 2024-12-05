package binder

type RFQCreateRequest struct {
	RfqId     string           `json:"id_rfq"`
	OrderDate string           `json:"order_date"`
	Status    string           `json:"status"`
	VendorId  string           `json:"id_vendor"`
	Products  []ProductRequest `json:"products"`
}

type RFQUpdateRequest struct {
	RfqId     string           `json:"id_rfq"`
	ProductId string           `json:"id_product"`
	OrderDate string           `json:"order_date"`
	VendorId  string           `json:"id_vendor"`
	Status    string           `json:"status"`
	Products  []ProductRequest `json:"products"`
}

type ProductRequest struct {
	RfqsProductId string `json:"id_rfqproduct"`
	ProductId     string `json:"id_product"`
	VendorId      string `json:"id_vendor"`
	ProductName   string `json:"productname"`
	Quantity      string `json:"quantity"`
	UnitPrice     string `json:"unitprice"`
	Tax           string `json:"tax"`
	Subtotal      string `json:"subtotal"`
}

type UpdateRfqStatusRequest struct {
	RfqId string `json:"id_rfq" validate:"required"`
}

type RFQDeleteRequest struct {
	RfqId string `param:"id_rfq" validate:"required"`
}
