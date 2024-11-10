package binder

type VendorCreateRequest struct {
	Vendorname string `json:"vendorname" validate:"required"`
	Addressone string `json:"addressone" validate:"required"`
	Addresstwo string `json:"addresstwo" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Website    string `json:"website" validate:"required"`
}

type VendorUpdateRequest struct {
	VendorId   string `param:"id_vendor" validate:"required"`
	Vendorname string `json:"vendorname" validate:"required"`
	Addressone string `json:"addressone" validate:"required"`
	Addresstwo string `json:"addresstwo" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Website    string `json:"website" validate:"required"`
}

type VendorDeleteRequest struct {
	VendorId string `param:"id_vendor" validate:"required"`
}
