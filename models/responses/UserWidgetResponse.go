package responses

import (
	"encoding/json"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserWidgetData struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	SchemaID  uuid.UUID      `json:"schema_id"`
	Config    map[string]any `json:"config"`
	PositionX int32          `json:"position_x"`
	PositionY int32          `json:"position_y"`
	Width     int32          `json:"width"`
	Height    int32          `json:"height"`
	ZIndex    int32          `json:"z_index"`
	IsVisible bool           `json:"is_visible"`
}

type UserWidgetResponse struct {
	Data    *UserWidgetData `json:"data"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
}

func (u *UserWidgetData) UserWidgetFromData(data *repository.UserWidget) error {
	u.ID = data.ID
	u.UserID = data.UserID
	u.SchemaID = data.SchemaID
	u.PositionX = data.PositionX
	u.PositionY = data.PositionY
	u.Width = data.Width
	u.Height = data.Height
	u.ZIndex = data.ZIndex
	u.IsVisible = data.IsVisible

	if err := json.Unmarshal(data.Config, &u.Config); err != nil {
		return err
	} else {
		return nil
	}
}

func NewUserWidgetResponse() *UserWidgetResponse {
	return &UserWidgetResponse{Success: false, Message: ""}
}

func (u *UserWidgetResponse) Successful(
	ctx echo.Context,
	data *repository.UserWidget,
) error {
	u.Data = &UserWidgetData{}
	u.Data.UserWidgetFromData(data)
	u.Success = true
	u.Message = "OK"
	return ctx.JSON(200, u)
}

func (u *UserWidgetResponse) Fail(
	ctx echo.Context,
	code int,
	err error,
) error  {
	u.Message = err.Error()
	return ctx.JSON(code, u)
}
