package responses

type FoldersResponse struct {
	Data []FolderData `json:"data"`
	Success bool `json:"success"`
	Message string `json:"message"`
}

func NewFoldersResponse() *FolderResponse {
	return &FolderResponse{Success: false, Message: ""}
}

