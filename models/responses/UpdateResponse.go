package responses

type UpdateResponse struct {
	Data    *UserData `json:"data"`
	JWT     string   `json:"jwt"`
	Success bool     `json:"success"`
	Message string   `json:"message"`
} // @name UpdateUserResponse
