package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/google/uuid"
)

type ApiKeyService struct {
	ctx               context.Context
	config            *configuration.Configuration
	redisService      IRedisService
	defaultExpiration time.Duration
}

func CreateApiKeyService(
	ctx context.Context,
	config *configuration.Configuration,
	redisService IRedisService,
) *ApiKeyService {
	// config.ApiKey.DefaultExpiration = 2592000000000000 as 720 hours
	// we need to get this from the config as a duration that can be parsed
    defaultExpiration, err := time.ParseDuration(config.ApiKey.DefaultExpiration)
	if err != nil {
		defaultExpiration = 720 * time.Hour
	}
	return &ApiKeyService{
		ctx:               ctx,
		config:            config,
		redisService:      redisService,
		defaultExpiration: defaultExpiration,
	}
}

func redisKeyForUser(userID uuid.UUID) string {
	return fmt.Sprintf("apikeys:%s", userID.String())
}

func generateRandomKey() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func hashKey(rawKey string, userID uuid.UUID) string {
	h := sha256.New()
	h.Write([]byte(rawKey + userID.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *ApiKeyService) Create(userID uuid.UUID, params CreateApiKeyParams) (*ApiKeyDTO, error) {
	rawKey, err := generateRandomKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	hashed := hashKey(rawKey, userID)

	expiration := params.Expiration
	if expiration == nil && s.defaultExpiration > 0 {
		t := time.Now().Add(s.defaultExpiration)
		expiration = &t
	}

	apiKey := &ApiKey{
		KeyID:       uuid.New(),
		UserID:      userID,
		HashedKey:   hashed,
		Name:        params.Name,
		NotBefore:   params.NotBefore,
		Expiration:  expiration,
		WriteAccess: params.WriteAccess,
		ReadAccess:  params.ReadAccess,
		CreatedAt:   time.Now(),
	}

	jsonStr, err := apiKey.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize API key: %w", err)
	}

	if err := s.redisService.LPush(redisKeyForUser(userID), jsonStr); err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}

	compositeKey := fmt.Sprintf("%s.%s", userID.String(), rawKey)

	return &ApiKeyDTO{
		KeyID:       apiKey.KeyID,
		Name:        apiKey.Name,
		NotBefore:   apiKey.NotBefore,
		Expiration:  apiKey.Expiration,
		WriteAccess: apiKey.WriteAccess,
		ReadAccess:  apiKey.ReadAccess,
		CreatedAt:   apiKey.CreatedAt,
		RawKey:      compositeKey,
	}, nil
}

func (s *ApiKeyService) List(userID uuid.UUID) ([]ApiKeyDTO, error) {
	entries, err := s.redisService.LRange(redisKeyForUser(userID), 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	dtos := make([]ApiKeyDTO, 0, len(entries))
	for _, entry := range entries {
		key, err := ApiKeyFromJSON(entry)
		if err != nil {
			continue
		}
		dtos = append(dtos, ApiKeyDTO{
			KeyID:       key.KeyID,
			Name:        key.Name,
			NotBefore:   key.NotBefore,
			Expiration:  key.Expiration,
			WriteAccess: key.WriteAccess,
			ReadAccess:  key.ReadAccess,
			CreatedAt:   key.CreatedAt,
		})
	}

	return dtos, nil
}

func (s *ApiKeyService) Delete(userID uuid.UUID, keyID uuid.UUID) error {
	entries, err := s.redisService.LRange(redisKeyForUser(userID), 0, -1)
	if err != nil {
		return fmt.Errorf("failed to read API keys: %w", err)
	}

	for _, entry := range entries {
		key, err := ApiKeyFromJSON(entry)
		if err != nil {
			continue
		}
		if key.KeyID == keyID {
			if err := s.redisService.LRem(redisKeyForUser(userID), 1, entry); err != nil {
				return fmt.Errorf("failed to delete API key: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("API key not found")
}

func (s *ApiKeyService) Validate(compositeKey string) (*ApiKey, error) {
	parts := strings.SplitN(compositeKey, ".", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid API key format")
	}

	userID, err := uuid.Parse(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid API key format: invalid user ID")
	}
	rawKey := parts[1]

	hashed := hashKey(rawKey, userID)

	entries, err := s.redisService.LRange(redisKeyForUser(userID), 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	now := time.Now()
	for _, entry := range entries {
		key, err := ApiKeyFromJSON(entry)
		if err != nil {
			continue
		}
		if key.HashedKey != hashed {
			continue
		}

		if key.Expiration != nil && now.After(*key.Expiration) {
			return nil, fmt.Errorf("API key has expired")
		}
		if key.NotBefore != nil && now.Before(*key.NotBefore) {
			return nil, fmt.Errorf("API key is not yet valid")
		}

		return key, nil
	}

	return nil, fmt.Errorf("invalid API key")
}
