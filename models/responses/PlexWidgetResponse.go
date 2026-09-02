package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type PlexMediaItemData struct {
	Title       string `json:"title"`
	Thumb       string `json:"thumb"`
	AddedAt     int64  `json:"added_at"`
	Type        string `json:"type"`
	Year        int    `json:"year"`
	ParentTitle string `json:"parent_title"`
} // @name PlexMediaItemData

type PlexRecentlyAddedResponse struct {
	Data    []PlexMediaItemData `json:"data"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
} // @name PlexRecentlyAddedResponse

func NewPlexRecentlyAddedResponse() *PlexRecentlyAddedResponse {
	return &PlexRecentlyAddedResponse{
		Data:    []PlexMediaItemData{},
		Success: true,
		Message: "Ok",
	}
}

func (r *PlexRecentlyAddedResponse) Fail(ctx echo.Context, code int, err error) error {
	r.Success = false
	r.Message = err.Error()
	return ctx.JSON(code, r)
}

func (r *PlexRecentlyAddedResponse) Successful(ctx echo.Context, items []clients.PlexMetadata) error {
	r.Success = true
	r.Data = make([]PlexMediaItemData, len(items))
	for i, item := range items {
		r.Data[i] = PlexMediaItemData{
			Title:       item.Title,
			Thumb:       item.Thumb,
			AddedAt:     item.AddedAt,
			Type:        item.Type,
			Year:        item.Year,
			ParentTitle: item.ParentTitle,
		}
	}
	return ctx.JSON(200, r)
}
