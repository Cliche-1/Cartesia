package middleware

import (
	"Gin/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	UserIDKey          = "user_id"
)

type AuthMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthMiddleware(jwtService *services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth verifica que el token JWT sea válido
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Se requiere autenticación",
			})
			return
		}

		// Extraer el token del header "Bearer <token>"
		tokenParts := strings.Split(header, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			return
		}

		// Validar el token
		userID, err := m.jwtService.ValidateToken(tokenParts[1])
		if err != nil {
			status := http.StatusUnauthorized
			message := "Token inválido"

			if err == services.ErrExpiredToken {
				message = "Token expirado"
			}

			c.AbortWithStatusJSON(status, gin.H{
				"error": message,
			})
			return
		}

		// Guardar el user_id en el contexto
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// GetUserID obtiene el ID del usuario del contexto
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}