package requests

import "time"

type CreateApiKeyRequest struct {
	Name        string     `json:"name"`
	NotBefore   *time.Time `json:"not_before,omitempty"`
	Expiration  *time.Time `json:"expiration,omitempty"`
	WriteAccess bool       `json:"write_access"`
	ReadAccess  bool       `json:"read_access"`
} // @name CreateApiKeyRequest
