package service

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type BillrfqService interface {
	CreateBill(mo *entity.Billrfq) (*entity.Billrfq, error)
}

type billrfqService struct {
	billrfqRepository repository.BillrfqRepository
}

func NewBillrfqService(billrfqRepository repository.BillrfqRepository) *billrfqService {
	return &billrfqService{
		billrfqRepository: billrfqRepository,
	}
}

func (s *billrfqService) CreateBill(mo *entity.Billrfq) (*entity.Billrfq, error) {

	lastId, err := s.billrfqRepository.GetLastMo()
	if err != nil {
		return nil, err
	}

	newMo := entity.NewBillrfq(lastId, mo.VendorId, mo.Bill_date, mo.Payment)

	savedMo, err := s.billrfqRepository.CreateMo(newMo)
	if err != nil {
		return nil, err
	}

	return savedMo, nil
}
