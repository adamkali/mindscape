package services

import (
	"context"
	"testing"
	"time"

	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// <method var=services.AuthService>
// <fixtures/>

// NewAuthTestConfig provides a base configuration for AuthService tests
func NewAuthTestConfig() *configuration.Configuration {
	config := &configuration.Configuration{}
	config.Server.JWT = "test-jwt-secret"
	return config
}

// WithAccessTTL modifies the config with a custom access token TTL
func WithAccessTTL(config *configuration.Configuration, ttl string) *configuration.Configuration {
	config.Server.AccessTokenTTL = ttl
	return config
}

// NewAuthTestUser provides a base user for AuthService tests
func NewAuthTestUser() *repository.User {
	return &repository.User{
		ID:       uuid.New(),
		Username: "authtestuser",
		Email:    "authtest@example.com",
		Admin:    false,
	}
}

// NewAuthTestService creates an AuthService without a database connection.
// Only stateless methods (CheckToken, MintAccessToken) may be exercised.
func NewAuthTestService(config *configuration.Configuration) *AuthService {
	return CreateAuthService(context.Background(), nil, config)
}

// signExpiredToken signs a token whose expiry is already in the past
func signExpiredToken(t *testing.T, config *configuration.Configuration, user *repository.User) string {
	claims := jwtFromUser(user, -time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Server.JWT))
	assert.NoError(t, err)
	return signed
}

// <runners/>

// Run_CheckToken_Valid verifies a freshly minted access token passes
func Run_CheckToken_Valid(t *testing.T) {
	config := NewAuthTestConfig()
	service := NewAuthTestService(config)

	access, err := service.MintAccessToken(NewAuthTestUser())
	assert.NoError(t, err)
	assert.NotEmpty(t, access)

	assert.NoError(t, service.CheckToken(access))
}

// Run_CheckToken_Expired verifies an expired access token is rejected
func Run_CheckToken_Expired(t *testing.T) {
	config := NewAuthTestConfig()
	service := NewAuthTestService(config)

	expired := signExpiredToken(t, config, NewAuthTestUser())
	assert.Error(t, service.CheckToken(expired))
}

// Run_CheckToken_WrongSignature verifies a token signed with a different key is rejected
func Run_CheckToken_WrongSignature(t *testing.T) {
	otherConfig := NewAuthTestConfig()
	otherConfig.Server.JWT = "some-other-secret"
	otherService := NewAuthTestService(otherConfig)

	access, err := otherService.MintAccessToken(NewAuthTestUser())
	assert.NoError(t, err)

	service := NewAuthTestService(NewAuthTestConfig())
	assert.Error(t, service.CheckToken(access))
}

// Run_CheckToken_ApiKeyBypass verifies the synthetic apikey token passes
func Run_CheckToken_ApiKeyBypass(t *testing.T) {
	service := NewAuthTestService(NewAuthTestConfig())
	assert.NoError(t, service.CheckToken("apikey"))
}

// Run_CheckToken_Garbage verifies a non-JWT string is rejected
func Run_CheckToken_Garbage(t *testing.T) {
	service := NewAuthTestService(NewAuthTestConfig())
	assert.Error(t, service.CheckToken("not-a-jwt"))
}

// Run_MintAccessToken_TTL verifies the configured TTL drives token expiry
func Run_MintAccessToken_TTL(t *testing.T) {
	config := WithAccessTTL(NewAuthTestConfig(), "5m")
	service := NewAuthTestService(config)

	access, err := service.MintAccessToken(NewAuthTestUser())
	assert.NoError(t, err)

	claims := &CustomJwt{}
	_, err = jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Server.JWT), nil
	})
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), claims.ExpiresAt.Time, 10*time.Second)
}

// Run_AccessTokenDuration_Defaults verifies fallback when TTL is unset or invalid
func Run_AccessTokenDuration_Defaults(t *testing.T) {
	config := NewAuthTestConfig()
	assert.Equal(t, configuration.DefaultAccessTokenTTL, config.AccessTokenDuration())
	assert.Equal(t, configuration.DefaultRefreshTokenTTL, config.RefreshTokenDuration())

	config.Server.AccessTokenTTL = "garbage"
	config.Server.RefreshTokenTTL = "-1h"
	assert.Equal(t, configuration.DefaultAccessTokenTTL, config.AccessTokenDuration())
	assert.Equal(t, configuration.DefaultRefreshTokenTTL, config.RefreshTokenDuration())

	config.Server.AccessTokenTTL = "30m"
	assert.Equal(t, 30*time.Minute, config.AccessTokenDuration())
}

// Run_GenerateRefreshToken_UniqueAndHashed verifies refresh tokens are
// high-entropy, unique per call, and deterministically hashed
func Run_GenerateRefreshToken_UniqueAndHashed(t *testing.T) {
	rawA, hashA, err := generateRefreshToken()
	assert.NoError(t, err)
	rawB, hashB, err := generateRefreshToken()
	assert.NoError(t, err)

	assert.NotEqual(t, rawA, rawB)
	assert.NotEqual(t, hashA, hashB)
	// the raw token never equals its stored hash
	assert.NotEqual(t, rawA, hashA)
	// hashing is deterministic so lookups by hash work
	assert.Equal(t, hashA, hashRefreshToken(rawA))
	// 32 bytes base64url-encoded => 43 chars
	assert.Len(t, rawA, 43)
}

// <tests>
// <map/>

// AuthServiceTestMap defines all AuthService test cases
var AuthServiceTestMap = map[string]func(*testing.T){
	"CheckToken_Valid":                  Run_CheckToken_Valid,
	"CheckToken_Expired":                Run_CheckToken_Expired,
	"CheckToken_WrongSignature":         Run_CheckToken_WrongSignature,
	"CheckToken_ApiKeyBypass":           Run_CheckToken_ApiKeyBypass,
	"CheckToken_Garbage":                Run_CheckToken_Garbage,
	"MintAccessToken_TTL":               Run_MintAccessToken_TTL,
	"AccessTokenDuration_Defaults":      Run_AccessTokenDuration_Defaults,
	"GenerateRefreshToken_UniqueHashed": Run_GenerateRefreshToken_UniqueAndHashed,
}

// <hook/>

// Test_AuthService tests all AuthService scenarios
func Test_AuthService(t *testing.T) {
	for name, testFunc := range AuthServiceTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>
