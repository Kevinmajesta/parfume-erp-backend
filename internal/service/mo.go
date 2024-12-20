package service

import (
	"bytes"
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type MoService interface {
	CreateMo(mo *entity.Mos) (*entity.Mos, error)
	UpdateMoStatus(moId string) (*entity.Mos, error)
	FindAllMos(page int) ([]entity.Mos, error)
	GetMoByID(MoId string) (*entity.Mos, error)
	DeleteMo(MoId string) (bool, error)
	GenerateMOPDF(mo *entity.Mos) ([]byte, error)
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

func (s *moService) UpdateMoStatus(moId string) (*entity.Mos, error) {
	// Fetch the existing Manufacture Order
	mo, err := s.moRepository.FindMoByID(moId)
	if err != nil {
		return nil, errors.New("manufacture order not found")
	}

	// Cycle through statuses
	switch mo.Status {
	case "draft":
		mo.Status = "confirmed"
	case "confirmed":
		mo.Status = "on progress"
	case "on progress":
		mo.Status = "done"
	default:
		return nil, errors.New("invalid status transition")
	}

	// Save the updated Manufacture Order
	updatedMo, err := s.moRepository.UpdateMoStatus(mo)
	if err != nil {
		return nil, errors.New("failed to update manufacture order status")
	}

	return updatedMo, nil
}

func (s *moService) FindAllMos(page int) ([]entity.Mos, error) {
	return s.moRepository.FindAllMos(page)
}

func (s *moService) GetMoByID(MoId string) (*entity.Mos, error) {

	material, err := s.moRepository.FindMoByID(MoId)
	if err != nil {
		return nil, err
	}
	return material, nil
}

func (s *moService) DeleteMo(MoId string) (bool, error) {
	material, err := s.moRepository.FindMoByID(MoId)
	if err != nil {
		return false, err
	}

	return s.moRepository.DeleteMo(material)
}

func (s *moService) GenerateMOPDF(mo *entity.Mos) ([]byte, error) {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 12, "Manufacturing Order (MO)", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Reset font to normal
	pdf.SetFont("Arial", "", 12)

	// Add MO details
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "MO ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, mo.MoId)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Product ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, mo.ProductId)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "BOM ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, mo.BomId)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Quantity to Produce:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, mo.Qtytoproduce)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Status:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, mo.Status)
	pdf.Ln(10)

	// Output the PDF to a buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	// Return the PDF as a byte slice
	return buf.Bytes(), nil
}
