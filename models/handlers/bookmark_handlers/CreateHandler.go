package bookmark_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/services"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateHandler struct {
	UserID          uuid.UUID
	Params          *repository.CreateBookmarkParams
	Bookmark        *repository.Bookmark
	Context         echo.Context
	Error           error
	Code            int
	Locked          bool
}


func NewCreateHandler(context echo.Context) *CreateHandler {
	return &CreateHandler{
		Context: context,
		Locked:  false,
		Code:    200,
	}
}
func (h *CreateHandler) Lock(code int) *CreateHandler {
	h.Locked = true
	h.Code = code
	return h
}
	
func (h *CreateHandler) Handle(fun any) *CreateHandler {
	var code int
	if !h.Locked {
		switch handle := fun.(type) {
		case func(token string) error: // this is to the jwt token
			code = 401
			h.Error = handle(h.Context.Get("user").(*jwt.Token).Raw)
			if h.Error != nil {
				return h.Lock(code)
			}
			jwt_token := h.Context.Get("user").(*jwt.Token)
			claims := jwt_token.Claims.(*services.CustomJwt)
			h.UserID = claims.UserId
			break
		case func(echo.Context) (*repository.CreateBookmarkParams, error):
			code = 400
			h.Params, h.Error = handle(h.Context)
			break
		case func(*repository.CreateBookmarkParams) (*repository.Bookmark, error):
			code = 500
			h.Bookmark, h.Error = handle(h.Params)
			break
		default:
			code = 600
			h.Error = echo.NewHTTPError(
				code,
				fmt.Sprintf("Type assertion failed for type: %T\n", fun),
			)
		}
		if h.Error != nil {
			return h.Lock(code)
		}
	}
	return h
}

func (h *CreateHandler) JSON() error {
	var code int
	var message string
	if h.Locked && h.Error != nil {
		code = h.Code
		if code == 600 {
			message = "Misaligend handler on the server: " + h.Error.Error()
		} else {
			message = h.Error.Error()
		}
	} else if code == 200 {
		message = "OK"
	}
	return h.Context.JSON(code, responses.NewBookmarkResponse(
		h.Bookmark,
		!h.Locked,
		message,
	))
}


