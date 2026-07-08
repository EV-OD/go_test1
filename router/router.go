package router

import (
	"myapp/config"
	"myapp/internal/user"

	_ "myapp/docs"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// @title Simple Api
// @version 1.0
// @description This is a simple API server for user management and transactions.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(cfg *config.Config) *echo.Echo {

	db, err := config.ConnectDB(cfg)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")
	api.GET("/swagger/*", echoSwagger.WrapHandler)

	user.RegisterRoutes(api, db, cfg)

	return e

}
