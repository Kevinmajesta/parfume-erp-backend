package binder

type CostumerCreateRequest struct {
	Costumername string `json:"costumername" validate:"required"`
	Addressone   string `json:"addressone" validate:"required"`
	Addresstwo   string `json:"addresstwo" `
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Status       string `json:"status"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
	City         string `json:"city"`
}

type CostumerUpdateRequest struct {
	CostumerId   string `param:"id_costumer" validate:"required"`
	Costumername string `json:"costumername" validate:"required"`
	Addressone   string `json:"addressone" validate:"required"`
	Addresstwo   string `json:"addresstwo" `
	Phone        string `json:"phone" `
	Email        string `json:"email" `
	Status       string `json:"status" `
	State        string `json:"state" `
	Zip          string `json:"zip" `
	Country      string `json:"country" `
	City         string `json:"city" `
}

type CostumerDeleteRequest struct {
	CostumerId string `param:"id_costumer" validate:"required"`
}
