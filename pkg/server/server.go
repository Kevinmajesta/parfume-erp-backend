package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
}

func NewServer(serverName string, publicRoutes, privateRoutes []*route.Route) *Server {
	e := echo.New()
	e.Use(middleware.CORS())

	e.Static("/assets", "assets")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:8000/"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Hello, World!", nil))
	})

	v1 := e.Group(fmt.Sprintf("/%s/api/v1", serverName))

	for _, public := range publicRoutes {
		v1.Add(public.Method, public.Path, public.Handler)
	}

	for _, private := range privateRoutes {
		// v1.Add(private.Method, private.Path, private.Handler, JWTProtected(cfg.JWT.SecretKey), SessionProtected())
		v1.Add(private.Method, private.Path, private.Handler)
	}
	return &Server{e}
}

func (s *Server) Run() {
	runServer(s)
	gracefulShutdown(s)
}

func runServer(srv *Server) {
	go func() {
		err := srv.Start(":3000")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func gracefulShutdown(srv *Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Fatal("Server Shutdown:", err)
	}
}


func contains(slice []string, s string) bool {
	for _, value := range slice {
		if value == s {
			return true
		}
	}
	return false
}
