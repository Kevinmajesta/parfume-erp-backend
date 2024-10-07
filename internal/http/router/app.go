package router

import (
	"net/http"

	"github.com/Kevinmajesta/webPemancingan/internal/http/handler"
	"github.com/Kevinmajesta/webPemancingan/pkg/route"
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
	schedulesHandler handler.SchedulesHandler) []*route.Route {
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
	}
}