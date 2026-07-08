package user

import (
	"myapp/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(grp *echo.Group, db *gorm.DB, cfg *config.Config) {

	userRepo := NewPostgresRepository(db)
	userService := NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	grp.POST("/register", userHandler.Register)
	grp.POST("/login", userHandler.Login)

	jwtConfig := echojwt.Config{
		SigningKey: []byte(cfg.JWT.Secret),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(CustomClaims)
		},

		SuccessHandler: userHandler.UserSuccessHandler,
	}

	protectedAPI := grp.Group("")
	protectedAPI.Use(echojwt.WithConfig(jwtConfig))

	protectedAPI.GET("/me/", userHandler.GetProfile)
	protectedAPI.POST("/loadbalance", userHandler.PostLoadBalanceHandler)

	roleBasedApi := protectedAPI.Group("")
	roleBasedApi.Use(userHandler.EnsureRole("editor"))
	roleBasedApi.POST("/send", userHandler.PostSendMoneyHandler)
}
