package entity

import (
	"fmt"
)

type Costumers struct {
	CostumerId   string `json:"id_costumer" gorm:"column:id_costumer;primaryKey"`
	Costumername string `json:"costumername"`
	Addressone   string `json:"addressone"`
	Addresstwo   string `json:"addresstwo"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Status       string `json:"status"`
	State        string `json:"state"`
	City         string `json:"city"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
	Auditable
}

func generateCostumerId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "CSR-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("CSR-%05d", newNumber)
}

func NewCostumer(lastId, costumername, addressone,
	addresstwo, phone, email, status, state, city, zip, country string) *Costumers {
	return &Costumers{
		CostumerId:   generateCostumerId(lastId),
		Costumername: costumername,
		Addressone:   addressone,
		Addresstwo:   addresstwo,
		Phone:        phone,
		Email:        email,
		Status:       status,
		State:        state,
		City:         city,
		Zip:          zip,
		Country:      country,
		Auditable:    NewAuditable(),
	}
}

func UpdateCostumer(id_costumer, costumername, addressone,
	addresstwo, phone, email, status, state, city, zip, country string) *Costumers {
	return &Costumers{
		CostumerId:   id_costumer,
		Costumername: costumername,
		Addressone:   addressone,
		Addresstwo:   addresstwo,
		Phone:        phone,
		Email:        email,
		Status:       status,
		State:        state,
		City:         city,
		Zip:          zip,
		Country:      country,
		Auditable:    UpdateAuditable(),
	}
}
