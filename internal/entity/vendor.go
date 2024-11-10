package entity

import (
	"fmt"
)

type Vendors struct {
	VendorId   string `json:"id_vendor" gorm:"column:id_vendor;primaryKey"`
	Vendorname string `json:"vendorname"`
	Addressone string `json:"addressone"`
	Addresstwo string `json:"addresstwo"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Website    string `json:"website"`
	Auditable
}

func generateVendorId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "VDR-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("VDR-%05d", newNumber)
}

func NewVendor(lastId, vendorname, addressone,
	addresstwo, phone, email, website string) *Vendors {
	return &Vendors{
		VendorId:   generateVendorId(lastId),
		Vendorname: vendorname,
		Addressone: addressone,
		Addresstwo: addresstwo,
		Phone:      phone,
		Email:      email,
		Website:    website,
		Auditable:  NewAuditable(),
	}
}

func UpdateVendor(id_vendor, vendorname, addressone,
	addresstwo, phone, email, website string) *Vendors {
	return &Vendors{
		VendorId:   id_vendor,
		Vendorname: vendorname,
		Addressone: addressone,
		Addresstwo: addresstwo,
		Phone:      phone,
		Email:      email,
		Website:    website,
		Auditable:  UpdateAuditable(),
	}
}
