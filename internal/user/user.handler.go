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
	EnsureRole(requiredRole string) echo.MiddlewareFunc
}

func NewUserHandler(s Service) Handler {
	return &userHandler{service: s}
}

type userHandler struct {
	service Service
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account profile in the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      RegisterRequest true "User Registration Payload"
// @Success      201     {object}  AuthResponse
// @Failure      400     {object}  apperrors.ApiError "Invalid request / Invalid input data"
// @Router       /register [post]
func (h *userHandler) Register(c *echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidRequest.Response(c)
	}

	authResponse, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		return apperrors.ErrorResponse(c, err)
	}

	return c.JSON(201, authResponse)
}

// Login godoc
// @Summary      Login a user
// @Description  Authenticate user and return JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      LoginRequest true "User Login Payload"
// @Success      200     {object}  AuthResponse
// @Failure      400     {object}  apperrors.ApiError "Invalid request / Invalid email or password"
// @Router       /login [post]
func (h *userHandler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidRequest.Response(c)
	}

	authResponse, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return ErrInvalidCredentials.Response(c)
	}

	return c.JSON(200, authResponse)
}

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Return profile for the authenticated user
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200     {object}  UserDTO
// @Failure      401     {object}  apperrors.ApiError "Unauthorized"
// @Router       /profile [get]
func (h *userHandler) GetProfile(c *echo.Context) error {
	uCtx := c.Get("currentUser")
	if uCtx == nil {
		return apperrors.ErrUnauthorized.Response(c)
	}

	eUser, ok := uCtx.(*EUser)
	if !ok {
		return ErrContextType.Response(c)
	}

	return c.JSON(http.StatusOK, eUser.UserDTO)
}

func (h *userHandler) UserSuccessHandler(c *echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return apperrors.ErrInvalidToken.Response(c)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return apperrors.ErrInvalidClaims.Response(c)
	}

	currentUserID := claims.ID
	user, err := h.service.GetFullProfile(c.Request().Context(), currentUserID)
	if err != nil {
		return ErrUserNotFound.Response(c)
	}

	c.Set("currentUser", user)
	c.Set("tokenString", token.Raw)
	return nil
}

// PostLoadBalanceHandler godoc
// @Summary      Load balance for user
// @Description  Add funds to authenticated user's balance
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body      LoadBalanceRequest true "Load Balance Payload"
// @Success      200     {object}  apperrors.ApiSuccess
// @Failure      400     {object}  apperrors.ApiError "Invalid request"
// @Failure      401     {object}  apperrors.ApiError "Unauthorized"
// @Router       /balance/load [post]
func (h *userHandler) PostLoadBalanceHandler(c *echo.Context) error {
	uCtx := c.Get("currentUser")
	if uCtx == nil {
		return apperrors.ErrUnauthorized.Response(c)
	}

	var req LoadBalanceRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidRequest.Response(c)
	}

	if err := h.service.LoadBalanceForUser(c.Request().Context(), uCtx.(*EUser).ID, req.Amount); err != nil {
		return apperrors.ErrorResponse(c, err)
	}

	return apperrors.SuccessResponse(c, "Balance loaded successfully")
}

// PostSendMoneyHandler godoc
// @Summary      Send money to another account
// @Description  Transfer funds from authenticated user to a receiver account
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body      SendMoneyRequest true "Send Money Payload"
// @Success      200     {object}  apperrors.ApiSuccess
// @Failure      400     {object}  apperrors.ApiError "Invalid request / insufficient balance"
// @Failure      401     {object}  apperrors.ApiError "Unauthorized"
// @Router       /money/send [post]
func (h *userHandler) PostSendMoneyHandler(c *echo.Context) error {
	uCtx := c.Get("currentUser")
	if uCtx == nil {
		return apperrors.ErrUnauthorized.Response(c)
	}

	var req SendMoneyRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.ErrInvalidRequest.Response(c)
	}

	if err := h.service.SendMoney(c.Request().Context(), uCtx.(*EUser).ID, req.ReceiverAccountID, req.Amount); err != nil {
		return apperrors.ErrorResponse(c, err)
	}

	return apperrors.SuccessResponse(c, "Money sent successfully")
}

func (h *userHandler) EnsureRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			user, ok := c.Get("currentUser").(*EUser)
			if !ok || user == nil {
				return ErrUserNotAuthorized.Response(c)
			}

			if user.Role != requiredRole {
				return ErrInSufficientRole.Response(c)
			}

			return next(c)
		}
	}
}
