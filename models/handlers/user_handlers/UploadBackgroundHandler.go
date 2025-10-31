package user_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/requests"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UploadBackgroundHandler struct {
	ctx              echo.Context
	code             int
	err              error
	url              string
	ValidatorService services.ValidatorService
	UserService      services.IUserService
	AuthService      services.IAuthService
	RedisService     services.IRedisService
	MinioService     services.IMinioService
}

func NewUploadBackgroundHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	UserService services.IUserService,
	AuthService services.IAuthService,
	RedisService services.IRedisService,
	MinioService services.IMinioService,
) *UploadBackgroundHandler {
	return &UploadBackgroundHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		UserService:      UserService,
		AuthService:      AuthService,
		RedisService:     RedisService,
		MinioService:     MinioService,
	}
}

// Handle handles the request from the frontend to upload a background
// to the minio bucket.
//
// it first checks if the file is valid,
// makes updates the user in the database,
// uploads the file to minio,
// gets a presigned url from minio,
// sets the new bacground url in redis with the expiration time,
// as the same in minio,
// and returns the presigned url as a responses.StringResponse
func (h *UploadBackgroundHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	var request *requests.UploadUserBackgroundRequest
	if request, err = h.ValidatorService.ValidateBacgroundImageChange(h.ctx); err != nil {
		return handlers.Lock(h, 400, err)
	}
	var params *repository.UpdateUserBacgroundParams
	file, err := request.GetFile()
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	defer file.Close()
	params = &repository.UpdateUserBacgroundParams{
		ID:         jwt_token.Claims.(*services.CustomJwt).UserId,
		Background: &request.File.Filename,
	}
	var user *repository.User
	if user, err = h.UserService.UpdateUserBackgroundImage(params); err != nil {
		return handlers.Lock(h, 500, err)
	}
	fmt.Printf("[INFO] UploadBackgroundHandler.Handle{ user: %v }\n", user)
	h.err = h.MinioService.Upload(
		user.ID,
		*user.Background,
		file,
		request.File.Size,
		"background",
	)
	if h.err != nil {
		return handlers.Lock(h, 500, h.err)
	}
	return h
}

func (h *UploadBackgroundHandler) JSON() error {
	if h.err != nil {
		// we just want to return the error
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	} else {
		// we want to return the presigned url for use in the front end
		return responses.NewStringResponse().Successful(h.ctx, h.url)
	}
}

func (h *UploadBackgroundHandler) Data() any {
	return h.url
}

func (h *UploadBackgroundHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h

}

func (h *UploadBackgroundHandler) Code() int {
	return h.code
}

func (h *UploadBackgroundHandler) Lock(code int, err error) handlers.IHandler {
	return handlers.Lock(h, code, err)
}

func (h *UploadBackgroundHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
func (h *UploadBackgroundHandler) Error() error {
	return h.err
}
