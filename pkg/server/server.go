package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kevinmajesta/webPemancingan/pkg/response"
	"github.com/Kevinmajesta/webPemancingan/pkg/route"
	"github.com/Kevinmajesta/webPemancingan/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
}

func NewServer(serverName string, publicRoutes, privateRoutes []*route.Route, secretKey string) *Server {
	e := echo.New()
	e.Use(middleware.CORS())

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
		v1.Add(private.Method, private.Path, private.Handler, JWTProtection(secretKey), JWTCheckRoles(private.Roles...))
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

func JWTProtection(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(secretKey),
	})
}

func JWTCheckRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusBadRequest, "you must login first"))
			}
			claims := user.Claims.(*token.JwtCustomClaims)
			// Check if the user has the required role
			if !contains(roles, claims.Role) {
				return c.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusBadRequest, "you don't have access"))
			}

			return next(c)
		}
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
