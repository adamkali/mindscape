package services

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	KeyID       uuid.UUID  `json:"key_id"`
	UserID      uuid.UUID  `json:"user_id"`
	HashedKey   string     `json:"hashed_key"`
	Name        string     `json:"name"`
	NotBefore   *time.Time `json:"not_before,omitempty"`
	Expiration  *time.Time `json:"expiration,omitempty"`
	WriteAccess bool       `json:"write_access"`
	ReadAccess  bool       `json:"read_access"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (a *ApiKey) ToJSON() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ApiKeyFromJSON(data string) (*ApiKey, error) {
	var key ApiKey
	if err := json.Unmarshal([]byte(data), &key); err != nil {
		return nil, err
	}
	return &key, nil
}
