package user_handlers

import (
	"time"

	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type DefaultBackgroundHandler struct {
	ctx              echo.Context
	code             int
	err              error
	url              string
	ValidatorService services.ValidatorService
	RedisService     services.IRedisService
	MinioService     services.IMinioService
}

func NewDefaultBackgroundHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	RedisService services.IRedisService,
	MinioService services.IMinioService,
) *DefaultBackgroundHandler {
	return &DefaultBackgroundHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		RedisService:     RedisService,
		MinioService:     MinioService,
	}
}

func (h *DefaultBackgroundHandler) Lock(code int, err error) *DefaultBackgroundHandler {
	h.code = code
	h.err = err
	return h
}

func (h *DefaultBackgroundHandler) Handle() handlers.IHandler{
	h.url, h.err = h.RedisService.Get("default")
	if h.err == nil {
		return h
	}
	h.url, h.err = h.MinioService.GetDefault()
	if h.err != nil {
		return handlers.Lock(h,500, h.err)
	}
	h.err = h.RedisService.SetWithExpiration(
		"default",
		h.url,
		time.Hour*24*7,
	)
	if h.err != nil {
		return handlers.Lock(h,500, h.err)
	}
	return h
}

func (h *DefaultBackgroundHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewStringResponse().Successful(h.ctx, h.url)
	}
}

func (h *DefaultBackgroundHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *DefaultBackgroundHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *DefaultBackgroundHandler) Code() int {
	return h.code
}

func (h *DefaultBackgroundHandler) Data() any {
	return h.url
}

func (h *DefaultBackgroundHandler) Error() error {
	return h.err
}

