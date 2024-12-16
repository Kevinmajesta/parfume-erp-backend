package binder

type BillRfqCreateRequest struct {
	VendorId  string `json:"vendorId" validate:"required"`
	Bill_date string `json:"bill_date" validate:"required"`
	Payment   string `json:"payment" validate:"required"`
}
