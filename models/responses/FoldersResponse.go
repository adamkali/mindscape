package responses

type FoldersResponse struct {
	Data    []FolderData `json:"data"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

func NewFoldersResponse(data []FolderData, success bool, message string) *FoldersResponse {
	return &FoldersResponse{Data: data, Success: success, Message: message}
}
