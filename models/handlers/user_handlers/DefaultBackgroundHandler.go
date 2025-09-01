package user_handlers

import (
	"errors"
	"fmt"
	"time"

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

func (h *DefaultBackgroundHandler) Handle() *DefaultBackgroundHandler {
	h.url, h.err = h.RedisService.Get("default")
	if h.err == nil {
		return h
	}
	h.url, h.err = h.MinioService.GetDefault()
	if h.err != nil {
		return h.Lock(500, h.err)
	}
	h.err = h.RedisService.SetWithExpiration(
		"default",
		h.url,
		time.Hour*24*7,
	)
	if h.err != nil {
		return h.Lock(500, h.err)
	}
	return h
}

func (h *DefaultBackgroundHandler) JSON() error {
	if h.err != nil {
		errorMessage := errors.New(fmt.Sprintf("%d Error: %s", h.code, h.err.Error()))
		return responses.NewStringResponse().Fail(h.ctx, h.code, errorMessage)
	} else {
		return responses.NewStringResponse().Successful(h.ctx, h.url)
	}
}
