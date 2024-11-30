package entity

import "fmt"

type Rfqs struct {
	RfqId     string        `json:"id_rfq" gorm:"column:id_rfq;primaryKey"`
	OrderDate string        `json:"order_date" gorm:"column:order_date"`
	VendorId  string        `json:"id_vendor" gorm:"column:id_vendor"`
	Status    string        `json:"status"`
	Products  []RfqsProduct `json:"products" gorm:"foreignKey:id_product"` // Gunakan foreignKey
	Auditable
}

type RfqsProduct struct {
	RfqsProductId string `json:"id_rfqproduct" gorm:"column:id_rfqproduct"`
	ProductId     string `json:"id_product" gorm:"column:id_product"`
	RfqId         string `json:"id_rfq" gorm:"column:id_rfq"`
	VendorId      string `json:"id_vendor" gorm:"column:id_vendor"` 
	ProductName   string `json:"productname" gorm:"column:productname"`
	Quantity      string `json:"quantity" gorm:"column:quantity"`
	UnitPrice     string `json:"unitprice" gorm:"column:unitprice"`
	Tax           string `json:"tax" gorm:"column:tax"`
	Subtotal      string `json:"subtotal" gorm:"column:subtotal"`
	Auditable
}

func generateRfqsId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "RFQ-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("RFQ-%05d", newNumber)
}

func NewRfqs(lastId, order_date, status, vendorId string) *Rfqs {
	return &Rfqs{
		RfqId:     generateRfqsId(lastId),
		OrderDate: order_date,
		VendorId:  vendorId,
		Status:    "RFQ",
		Auditable: NewAuditable(),
	}
}

func UpdateRfqs(id_rfq, order_date, newStatus, vendorId string, currentStatus string) *Rfqs {
	// Jika status baru kosong, gunakan status lama
	status := currentStatus
	if newStatus != "" {
		status = newStatus
	}

	return &Rfqs{
		RfqId:     id_rfq,
		VendorId:  vendorId,
		OrderDate: order_date,
		Status:    status,
		Auditable: UpdateAuditable(),
	}
}
