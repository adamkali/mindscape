package user_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type BackgroundChoicesHandler struct {
	ctx              echo.Context
	code             int
	err              error
	urls             []responses.BackgroundData
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

func (h *BackgroundChoicesHandler) Handle() handlers.IHandler {
	entites, err := h.MinioService.GetBackgroundChoices()
	if err != nil {
		return handlers.Lock(h,500, h.err)
	}
	for _, entity := range entites {
		h.urls = append(h.urls, responses.NewBackgroundsData(entity.Name, entity.URL))
	}
	return h
}

func (h *BackgroundChoicesHandler) JSON() error {
	if h.err != nil {
		return responses.NewBackgroundResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewBackgroundResponse().Successful(h.ctx, h.urls)
	}
}

func (h *BackgroundChoicesHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *BackgroundChoicesHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *BackgroundChoicesHandler) Code() int {
	return h.code
}

func (h *BackgroundChoicesHandler) Data() any {
	return h.urls
}

func (h *BackgroundChoicesHandler) Error() error {
	return h.err
}
