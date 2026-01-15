package requests

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type AddUserWidgetRequst struct {
	SchemaID  uuid.UUID `json:"schema_id"`
	Config    []byte    `json:"config"`
	PositionX int32     `json:"position_x"`
	PositionY int32     `json:"position_y"`
	Width     int32     `json:"width"`
	Height    int32     `json:"height"`
	ZIndex    int32     `json:"z_index"`
	IsVisible bool      `json:"is_visible"`
} // @name AddUserWidgetRequest

func (wr AddUserWidgetRequst) IntoRepositoryParams(userID uuid.UUID) *repository.CreateUserWidgetParams {
	return &repository.CreateUserWidgetParams{
		UserID:    userID,
		SchemaID:  wr.SchemaID,
		Config:    wr.Config,
		PositionX: wr.PositionX,
		PositionY: wr.PositionY,
		Width:     wr.Width,
		Height:    wr.Height,
		ZIndex:    wr.ZIndex,
		IsVisible: wr.IsVisible,
	}
}
