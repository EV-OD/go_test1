package router

import (
	"myapp/config"
	"myapp/internal/user"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewRouter(cfg *config.Config) *echo.Echo {

	db, err := config.ConnectDB(cfg.Database.URL)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	user.RegisterRoutes(api, db, cfg)

	return e

}
