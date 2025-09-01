package user_handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type BackgroundChoicesHandler struct {
	ctx              echo.Context
	code             int
	err              error
	urls             []string
	ValidatorService services.ValidatorService
	RedisService     services.IRedisService
	MinioService     services.IMinioService
}

func NewBackgroundChoicesHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	RedisService services.IRedisService,
	MinioService services.IMinioService,
) *BackgroundChoicesHandler {
	return &BackgroundChoicesHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		RedisService:     RedisService,
		MinioService:     MinioService,
	}
}

func (h *BackgroundChoicesHandler) Lock(code int, err error) *BackgroundChoicesHandler {
	h.code = code
	h.err = err
	return h
}

func (h *BackgroundChoicesHandler) Handle() *BackgroundChoicesHandler {
	var urlstring string
	urlstring, h.err = h.RedisService.Get("background_choices")
	if h.err == nil {
		h.urls = strings.Split(urlstring, ",")
		return h
	}
	h.urls , h.err = h.MinioService.GetBackgroundChoices()
	if h.err != nil {
		return h.Lock(500, h.err)
	}
	h.err = h.RedisService.SetWithExpiration(
		"background_choices",
		strings.Join(h.urls, ","),
		time.Hour*24*7,
	)
	if h.err != nil {
		return h.Lock(500, h.err)
	}
	return h
}

func (h *BackgroundChoicesHandler) JSON() error {
	if h.err != nil {
		errorMessage := errors.New(fmt.Sprintf("%d Error: %s", h.code, h.err.Error()))
		return responses.NewStringsResponse().Fail(h.ctx, h.code, errorMessage)
	} else {
		return responses.NewStringsResponse().Successful(h.ctx, h.urls)
	}
}
