package responses

import "github.com/labstack/echo/v4"

type (
	BackgroundsResponse struct {
		Data    []BackgroundData `json:"data"`
		Success bool             `json:"success"`
		Message string           `json:"message"`
	} // @name BackgroundResponse

	BackgroundData struct {
		Filename string `json:"filename"`
		Url      string `json:"url"`
	} // @name BackgroundData
)

func NewBackgroundResponse() *BackgroundsResponse { return &BackgroundsResponse{} }

func NewBackgroundsData(filename string, url string) BackgroundData {
	return BackgroundData{
		Filename: filename,
		Url:      url,
	}
}

func (b *BackgroundsResponse) Successful(ctx echo.Context, data []BackgroundData) error {
	return ctx.JSON(200, BackgroundsResponse{
		Data:    data,
		Success: true,
		Message: "OK",
	})
}

func (b *BackgroundsResponse) Fail(ctx echo.Context, code int, err error) error {
	return ctx.JSON(code, BackgroundsResponse{
		Data:    []BackgroundData{},
		Success: false,
		Message: err.Error(),
	})
}
