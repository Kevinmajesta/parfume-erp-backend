package entity

import "fmt"

type Quotations struct {
	QuotationsId string              `json:"id_quotation" gorm:"column:id_quotation;primaryKey"`
	OrderDate    string              `json:"order_date" gorm:"column:order_date"`
	CostumerId   string              `json:"id_costumer" gorm:"column:id_costumer"`
	Status       string              `json:"status"`
	Payment      string              `json:"payment"`
	Products     []QuotationsProduct `json:"products" gorm:"foreignKey:QuotationsId;references:QuotationsId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;joinTableForeignKey:QuotationsId;table:quotations_products"`
	Auditable
}

type QuotationsProduct struct {
	QuotationsProductId string `json:"id_quotationsproduct" gorm:"column:id_quotationsproduct"`
	ProductId           string `json:"id_product" gorm:"column:id_product"`
	QuotationsId        string `json:"id_quotation" gorm:"column:id_quotation"`
	CostumerId          string `json:"id_costumer" gorm:"column:id_costumer"`
	ProductName         string `json:"productname" gorm:"column:productname"`
	Quantity            string `json:"quantity" gorm:"column:quantity"`
	UnitPrice           string `json:"unitprice" gorm:"column:unitprice"`
	Tax                 string `json:"tax" gorm:"column:tax"`
	Subtotal            string `json:"subtotal" gorm:"column:subtotal"`
	Auditable
}

func generateQuoId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "QUO-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("QUO-%05d", newNumber)
}

func NewQuo(lastId, order_date, status, CostumerId, payment string) *Quotations {
	return &Quotations{
		QuotationsId: generateQuoId(lastId),
		OrderDate:    order_date,
		CostumerId:   CostumerId,
		Status:       "QUOTATION",
		Payment:      payment,
		Auditable:    NewAuditable(),
	}
}

func UpdateQuo(id_quotation, order_date, newStatus, CostumerId, payment string, currentStatus string) *Quotations {
	// Jika status baru kosong, gunakan status lama
	status := currentStatus
	if newStatus != "" {
		status = newStatus
	}

	return &Quotations{
		QuotationsId: id_quotation,
		CostumerId:   CostumerId,
		OrderDate:    order_date,
		Status:       status,
		Payment:      payment,
		Auditable:    UpdateAuditable(),
	}
}
