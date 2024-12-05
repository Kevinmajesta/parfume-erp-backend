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
	Status     string `json:"status"`
	State      string `json:"state"`
	City       string `json:"city"`
	Zip        string `json:"zip"`
	Country    string `json:"country"`
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
	addresstwo, phone, email, website, status, state, city, zip, country string) *Vendors {
	return &Vendors{
		VendorId:   generateVendorId(lastId),
		Vendorname: vendorname,
		Addressone: addressone,
		Addresstwo: addresstwo,
		Phone:      phone,
		Email:      email,
		Website:    website,
		Status:     status,
		State:      state,
		City:       city,
		Zip:        zip,
		Country:    country,
		Auditable:  NewAuditable(),
	}
}

func UpdateVendor(id_vendor, vendorname, addressone,
	addresstwo, phone, email, website, status, state, city, zip, country string) *Vendors {
	return &Vendors{
		VendorId:   id_vendor,
		Vendorname: vendorname,
		Addressone: addressone,
		Addresstwo: addresstwo,
		Phone:      phone,
		Email:      email,
		Website:    website,
		Status:     status,
		State:      state,
		City:       city,
		Zip:        zip,
		Country:    country,
		Auditable:  UpdateAuditable(),
	}
}
