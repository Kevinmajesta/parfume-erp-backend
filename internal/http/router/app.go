package router

import (
	"net/http"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/handler"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/route"
)

const (
	Admin = "admin"
	User  = "user"
)

var (
	allRoles  = []string{Admin, User}
	onlyAdmin = []string{Admin}
	onlyUser  = []string{User}
)

func PublicRoutes(userHandler handler.UserHandler, adminHandler handler.AdminHandler) []*route.Route {
	return []*route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.LoginUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: userHandler.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login/admin",
			Handler: adminHandler.LoginAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admins",
			Handler: adminHandler.CreateAdmin,
		},
	}
}

func PrivateRoutes(userHandler handler.UserHandler, suggestionHandler handler.SuggestionHandler, adminHandler handler.AdminHandler,
	schedulesHandler handler.SchedulesHandler, productHandler handler.ProductHandler, materialHandler handler.MaterialHandler,
	bomHandler handler.BOMHandler, moHandler handler.MoHandler, vendorHandler handler.VendorHandler, rfqHandler handler.RfqHandler,
	costumerHandler handler.CostumerHandler, quoHandler handler.QuoHandler, billrfqHandler handler.BillrfqHandler) []*route.Route {
	return []*route.Route{
		//user
		{
			Method:  http.MethodPut,
			Path:    "/users/:id_user",
			Handler: userHandler.UpdateUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/:id_user",
			Handler: userHandler.DeleteUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id_user",
			Handler: userHandler.GetUserProfile,
			Roles:   allRoles,
		},
		//suggestion
		{
			Method:  http.MethodPost,
			Path:    "/suggestions",
			Handler: suggestionHandler.CreateSuggestion,
			Roles:   allRoles,
		},
		//admin
		{
			Method:  http.MethodGet,
			Path:    "/allusers",
			Handler: adminHandler.FindAllUser,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPut,
			Path:    "/admins/:id_user",
			Handler: adminHandler.UpdateAdmin,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admins/:id_user",
			Handler: adminHandler.DeleteAdmin,
			Roles:   onlyAdmin,
		},
		//schedules
		{
			Method:  http.MethodPost,
			Path:    "/schedule",
			Handler: schedulesHandler.CreateSchedules,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/allschedule",
			Handler: schedulesHandler.FindAllSchedule,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPut,
			Path:    "/edit/schedule/:id_schedules",
			Handler: schedulesHandler.UpdateSchedule,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/delete/schedule/:id_schedules",
			Handler: schedulesHandler.DeleteSchedule,
			Roles:   onlyAdmin,
		},
		//product
		{
			Method:  http.MethodPost,
			Path:    "/products",
			Handler: productHandler.CreateProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPut,
			Path:    "/products/:id_product",
			Handler: productHandler.UpdateProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/products/:id_product",
			Handler: productHandler.DeleteProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product/all",
			Handler: productHandler.FindAllProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product/variants/all",
			Handler: productHandler.FindAllProductVariant,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product",
			Handler: productHandler.SearchProducts,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product/:id_product/pdf",
			Handler: productHandler.DownloadProductPDF,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product/:id_product",
			Handler: productHandler.GetProductProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product/pdf",
			Handler: productHandler.GenerateAllProductsPDFHandler,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/product/increase",
			Handler: productHandler.IncreaseProductQty,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/product/decrease",
			Handler: productHandler.DecreaseProductQty,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/product/:id_product/barcode",
			Handler: productHandler.GenerateBarcodePDFHandler,
			Roles:   allRoles,
		},
		//material
		{
			Method:  http.MethodPost,
			Path:    "/materials",
			Handler: materialHandler.CreateMaterial,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPut,
			Path:    "/materials/:id_material",
			Handler: materialHandler.UpdateMaterial,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/materials/:id_material",
			Handler: materialHandler.DeleteMaterial,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/material/all",
			Handler: materialHandler.FindAllMaterial,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/material",
			Handler: materialHandler.SearchMaterials,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/material/:id_material",
			Handler: materialHandler.GetMaterialProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/material/:id_material/pdf",
			Handler: materialHandler.DownloadMaterialPDF,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/material/pdf",
			Handler: materialHandler.GenerateAllMaterialsPDFHandler,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/material/:id_material/barcode",
			Handler: materialHandler.GenerateBarcodePDFHandler,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/material/reducemat",
			Handler: materialHandler.ReduceMaterialQty,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/material/increasemat",
			Handler: materialHandler.IncreaseMaterialQty,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/bom",
			Handler: bomHandler.CreateBOM,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/bom/all",
			Handler: bomHandler.FindAllBom,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/bom/:id_bom",
			Handler: bomHandler.DeleteBom,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/bom/edit/:id_bom",
			Handler: bomHandler.UpdateBOM,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/bom/:id_bom",
			Handler: bomHandler.GetBOMByID,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/bom/:id_bom/overview",
			Handler: bomHandler.GetBOMOverview,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/bom/:id_bom/overview/pdf",
			Handler: bomHandler.GetBOMPDF,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/mo",
			Handler: moHandler.CreateMo,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/mo/status/confirm",
			Handler: moHandler.UpdateMoStatus,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/mo/all",
			Handler: moHandler.FindAllMos,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/mo/:id_mo",
			Handler: moHandler.GetMoProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/mo/:id_mo",
			Handler: moHandler.DeleteMo,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/mo/:id_mo/pdf",
			Handler: moHandler.DownloadMOPDF,
			Roles:   allRoles,
		},
		//vendor
		{
			Method:  http.MethodPost,
			Path:    "/vendor",
			Handler: vendorHandler.CreateVendor,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/vendor/:id_vendor",
			Handler: vendorHandler.UpdateVendor,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/vendor/:id_vendor",
			Handler: vendorHandler.DeleteVendor,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/vendors",
			Handler: vendorHandler.FindAllVendor,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/vendor/:id_vendor",
			Handler: vendorHandler.GetVendorProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/vendor/:id_vendor/pdf",
			Handler: vendorHandler.DownloadVendorPDF,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/vendor/pdf",
			Handler: vendorHandler.DownloadAllVendorsPDF,
			Roles:   allRoles,
		},
		//RFQ
		{
			Method:  http.MethodPost,
			Path:    "/rfq",
			Handler: rfqHandler.CreateRfq,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/rfq",
			Handler: rfqHandler.UpdateRfq,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/rfq/:id_rfq",
			Handler: rfqHandler.UpdateRfqAll,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/rfq/status",
			Handler: rfqHandler.UpdateRfqStatus,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/rfq/all/rfq",
			Handler: rfqHandler.FindAllRfq,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/rfq/all/bill",
			Handler: rfqHandler.FindAllRfqBill,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/rfq/:id_rfq",
			Handler: rfqHandler.GetRfqOverview,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/rfq/email/:id_vendor",
			Handler: rfqHandler.GetVendorEmailById,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/rfq/:id_rfq",
			Handler: rfqHandler.DeleteRfq,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/rfq/:id_rfq/pdf",
			Handler: rfqHandler.HandleCreateRfqPDF,
			Roles:   allRoles,
		},
		//costumer
		{
			Method:  http.MethodPost,
			Path:    "/costumer",
			Handler: costumerHandler.CreateCostumer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/costumer/:id_costumer",
			Handler: costumerHandler.UpdateCostumer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/costumer/:id_costumer",
			Handler: costumerHandler.DeleteCostumer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/costumers",
			Handler: costumerHandler.FindAllCostumer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/costumer/:id_costumer",
			Handler: costumerHandler.GetCostumerProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/costumer/:id_costumer/pdf",
			Handler: costumerHandler.HandleCreateCostumerPDF,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/costumers/pdf",
			Handler: costumerHandler.HandleCreateCostumerPDFAll,
			Roles:   allRoles,
		},
		//quo
		{
			Method:  http.MethodPost,
			Path:    "/quotation",
			Handler: quoHandler.CreateRfq,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/quotation/:id_quotation",
			Handler: quoHandler.UpdateQuo,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/quotation/:id_quotation",
			Handler: quoHandler.DeleteRfq,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/quotation/all/quo",
			Handler: quoHandler.FindAllQuo,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/quotation/all/bill",
			Handler: quoHandler.FindAllQuoBill,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/quotation/status/:id_quotation",
			Handler: quoHandler.UpdateQuoStatus,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/quotation/overview/:id_quotation",
			Handler: quoHandler.GetQuoOverview,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/quotation/email/:id_costumer",
			Handler: quoHandler.GetCostumerEmailById,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/quotation/:id_quotation/pdf",
			Handler: quoHandler.HandleCreateQuoPDF,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/billrfq",
			Handler: billrfqHandler.CreateMo,
			Roles:   allRoles,
		},
	}
}
