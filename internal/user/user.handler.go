package user

import (
	"net/http"

	apperrors "myapp/internal/errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v5"
)

type Handler interface {
	Register(c *echo.Context) error
	Login(c *echo.Context) error
	GetProfile(c *echo.Context) error
	UserSuccessHandler(c *echo.Context) error
	PostLoadBalanceHandler(c *echo.Context) error
	PostSendMoneyHandler(c *echo.Context) error
}

func NewUserHandler(s Service) Handler {
	return &userHandler{service: s}
}

type userHandler struct {
	service Service
}

func (h *userHandler) Register(c *echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrInvalidRequest.Code), map[string]string{"error": apperrors.ErrInvalidRequest.Message})
	}

	authResponse, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(201, authResponse)
}

func (h *userHandler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrInvalidRequest.Code), map[string]string{"error": apperrors.ErrInvalidRequest.Message})
	}

	authResponse, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(apperrors.HTTPStatus(ErrInvalidCredentials.Code), map[string]string{"error": ErrInvalidCredentials.Message})
	}

	return c.JSON(200, authResponse)
}

func (h *userHandler) GetProfile(c *echo.Context) error {
	uCtx := c.Get("currentUser")
	if uCtx == nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrUnauthorized.Code), map[string]string{"error": apperrors.ErrUnauthorized.Message})
	}

	eUser, ok := uCtx.(*EUser)
	if !ok {
		return c.JSON(apperrors.HTTPStatus(ErrContextType.Code), map[string]string{"error": ErrContextType.Message})
	}

	return c.JSON(http.StatusOK, eUser.UserDTO)
}

func (h *userHandler) UserSuccessHandler(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrInvalidToken.Code), map[string]string{"error": apperrors.ErrInvalidToken.Message})
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrInvalidClaims.Code), map[string]string{"error": apperrors.ErrInvalidClaims.Message})
	}

	currentUserID := claims.ID
	user, err := h.service.GetFullProfile(c.Request().Context(), currentUserID)
	if err != nil {
		return c.JSON(apperrors.HTTPStatus(ErrUserNotFound.Code), map[string]string{"error": ErrUserNotFound.Message})
	}

	c.Set("currentUser", user)
	c.Set("tokenString", token.Raw)
	return nil
}

func (h *userHandler) PostLoadBalanceHandler(c *echo.Context) error {
	uCtx := c.Get("currentUser")
	if uCtx == nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrUnauthorized.Code), map[string]string{"error": apperrors.ErrUnauthorized.Message})
	}

	var req LoadBalanceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrInvalidRequest.Code), map[string]string{"error": apperrors.ErrInvalidRequest.Message})
	}

	if err := h.service.LoadBalanceForUser(c.Request().Context(), uCtx.(*EUser).ID, req.Amount); err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"message": "Balance loaded successfully"})
}

func (h *userHandler) PostSendMoneyHandler(c *echo.Context) error {
	uCtx := c.Get("currentUser")
	if uCtx == nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrUnauthorized.Code), map[string]string{"error": apperrors.ErrUnauthorized.Message})
	}

	var req SendMoneyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(apperrors.HTTPStatus(apperrors.ErrInvalidRequest.Code), map[string]string{"error": apperrors.ErrInvalidRequest.Message})
	}

	if err := h.service.SendMoney(c.Request().Context(), uCtx.(*EUser).ID, req.ReceiverAccountID, req.Amount); err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"message": "Money sent successfully"})
}
