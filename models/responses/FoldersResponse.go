package responses

import "github.com/labstack/echo/v4"

type FoldersResponse struct {
	Data    []FolderData `json:"data"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

func NewFoldersResponse() *FoldersResponse {
	return &FoldersResponse{Data: []FolderData{}, Success: false, Message: "Internal Server Error"}
}

func (f *FoldersResponse) Fail(ctx echo.Context, code int, err error) error {
	f.Success = false
	f.Message = err.Error()
	return ctx.JSON(code, f)
}

func (f *FoldersResponse) Successful(ctx echo.Context, data []FolderData) error {
	f.Success = true
	f.Message = "OK"
	for _, folder := range data {
		f.Data = append(f.Data, folder)
	}
	return ctx.JSON(200, f)
}
