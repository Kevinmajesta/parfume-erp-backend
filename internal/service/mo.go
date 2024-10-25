package service

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type MoService interface {
	CreateMo(mo *entity.Mos) (*entity.Mos, error)
}

type moService struct {
	moRepository repository.MoRepository
}

func NewMoService(moRepository repository.MoRepository) *moService {
	return &moService{
		moRepository: moRepository,
	}
}
func (s *moService) CreateMo(mo *entity.Mos) (*entity.Mos, error) {

	lastId, err := s.moRepository.GetLastMo()
	if err != nil {
		return nil, err
	}

	newMo := entity.NewMos(lastId, mo.ProductId, mo.BomId, mo.Qtytoproduce)

	savedMo, err := s.moRepository.CreateMo(newMo)
	if err != nil {
		return nil, err
	}

	return savedMo, nil
}

