package user_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


type UserBackgroundChoicesHandler struct {
	ctx              echo.Context
	code             int
	err              error
	urls             []responses.BackgroundData
	ValidatorService services.ValidatorService
	AuthService      services.IAuthService
	MinioService     services.IMinioService
}

func NewUserBackgroundChoicesHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	AuthService services.IAuthService,
	MinioService services.IMinioService,
) *UserBackgroundChoicesHandler {
	return &UserBackgroundChoicesHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		AuthService:      AuthService,
		MinioService:     MinioService,
	}
 }
func (h *UserBackgroundChoicesHandler) Handle() handlers.IHandler {
	fmt.Println("[INFO] UserBackgroundChoicesHandler.Handle")
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	err := h.AuthService.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	fmt.Printf("[INFO] UserBackgroundChoicesHandler.Handle{ userID: %v }\n", userID)
	entities, err := h.MinioService.GetUserBackgroundChoices(userID)
	if err != nil {
		return handlers.Lock(h, 500, err)
	}
	for _, entity := range entities {
		h.urls = append(h.urls, responses.NewBackgroundsData(entity.Name, entity.URL))
	}
	return h
}

func (h *UserBackgroundChoicesHandler) JSON() error {
	if h.err != nil {
		return responses.NewBackgroundResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewBackgroundResponse().Successful(h.ctx, h.urls)
	}
}

func (h *UserBackgroundChoicesHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *UserBackgroundChoicesHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *UserBackgroundChoicesHandler) Code() int {
	return h.code
}

func (h *UserBackgroundChoicesHandler) Data() any {
	return h.urls
}

func (h *UserBackgroundChoicesHandler) Error() error {
	return h.err
}
