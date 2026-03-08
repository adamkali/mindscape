package configs

import (
	"fmt"
	"strings"

	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ApiKeyMiddlewareConfig(apiKeyService services.IApiKeyService) middleware.KeyAuthConfig {
	return middleware.KeyAuthConfig{
		KeyLookup: "header:X-API-Key",
		Validator: func(key string, c echo.Context) (bool, error) {
			apiKey, err := apiKeyService.Validate(key)
			if err != nil {
				return false, err
			}

			method := c.Request().Method
			isWrite := method == "POST" || method == "PUT" || method == "DELETE" || method == "PATCH"

			if isWrite && !apiKey.WriteAccess {
				return false, fmt.Errorf("API key does not have write access")
			}
			if !isWrite && !apiKey.ReadAccess {
				return false, fmt.Errorf("API key does not have read access")
			}

			c.Set("apikey", apiKey)
			c.Set("apikey_user_id", apiKey.UserID)

			// Set a synthetic JWT token in context so existing handlers
			// that call c.Get("user").(*jwt.Token) can extract the userID.
			// The token Raw field is set to "apikey" so AuthService.CheckToken
			// is bypassed in API key-aware handlers.
			syntheticClaims := &services.CustomJwt{
				UserId: apiKey.UserID,
			}
			syntheticToken := &jwt.Token{
				Claims: syntheticClaims,
				Raw:    "apikey",
				Valid:  true,
			}
			c.Set("user", syntheticToken)

			return true, nil
		},
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "swagger")
		},
		ErrorHandler: func(err error, c echo.Context) error {
			if strings.Contains(err.Error(), "write access") {
				return c.JSON(403, map[string]string{"message": err.Error()})
			}
			return c.JSON(401, map[string]string{"message": err.Error()})
		},
	}
}
