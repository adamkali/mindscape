package services

import (
	"time"

	"github.com/google/uuid"
)

type CreateApiKeyParams struct {
	Name        string     `json:"name"`
	NotBefore   *time.Time `json:"not_before,omitempty"`
	Expiration  *time.Time `json:"expiration,omitempty"`
	WriteAccess bool       `json:"write_access"`
	ReadAccess  bool       `json:"read_access"`
}

type ApiKeyDTO struct {
	KeyID       uuid.UUID  `json:"key_id"`
	Name        string     `json:"name"`
	NotBefore   *time.Time `json:"not_before,omitempty"`
	Expiration  *time.Time `json:"expiration,omitempty"`
	WriteAccess bool       `json:"write_access"`
	ReadAccess  bool       `json:"read_access"`
	CreatedAt   time.Time  `json:"created_at"`
	RawKey      string     `json:"raw_key,omitempty"`
}

type IApiKeyService interface {
	Create(userID uuid.UUID, params CreateApiKeyParams) (*ApiKeyDTO, error)
	List(userID uuid.UUID) ([]ApiKeyDTO, error)
	Delete(userID uuid.UUID, keyID uuid.UUID) error
	Validate(compositeKey string) (*ApiKey, error)
}
