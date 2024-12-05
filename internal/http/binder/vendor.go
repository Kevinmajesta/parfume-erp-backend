package binder

type VendorCreateRequest struct {
	Vendorname string `json:"vendorname" validate:"required"`
	Addressone string `json:"addressone" validate:"required"`
	Addresstwo string `json:"addresstwo" `
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Website    string `json:"website" `
	Status     string `json:"status" validate:"required"`
	State 	string `json:"state" validate:"required"`
	Zip 	string `json:"zip" validate:"required"`
	Country 	string `json:"country" validate:"required"`
	City 	string `json:"city" validate:"required"`
}

type VendorUpdateRequest struct {
	VendorId   string `param:"id_vendor" validate:"required"`
	Vendorname string `json:"vendorname" validate:"required"`
	Addressone string `json:"addressone" validate:"required"`
	Addresstwo string `json:"addresstwo" `
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Website    string `json:"website" `
	Status     string `json:"status" validate:"required"`
	State 	string `json:"state" validate:"required"`
	Zip 	string `json:"zip" validate:"required"`
	Country 	string `json:"country" validate:"required"`
	City 	string `json:"city" validate:"required"`
}

type VendorDeleteRequest struct {
	VendorId string `param:"id_vendor" validate:"required"`
}
