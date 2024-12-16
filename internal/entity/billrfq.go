package entity

import "fmt"

type Billrfq struct {
	BillrfqId string `json:"id_bill" gorm:"column:id_bill;primaryKey"`
	VendorId  string `json:"id_vendor" gorm:"column:id_vendor"`
	Bill_date string `json:"bill_date"`
	Payment   string `json:"payment"`
	Auditable
}

func generateBillRfqId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "BRQ-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("BRQ-%05d", newNumber)
}

func NewBillrfq(lastId, id_vendor, bill_date, payment string) *Billrfq {
	return &Billrfq{
		BillrfqId: generateBillRfqId(lastId),
		VendorId:  id_vendor,
		Bill_date: bill_date,
		Payment:   payment,
		Auditable: NewAuditable(),
	}
}
