package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired es un middleware que verifica el token JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no proporcionado"})
			c.Abort()
			return
		}

		// El token debe estar en formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Validar el token
		token, err := validateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Si el token es válido, obtener los claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Guardar el ID del usuario en el contexto
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}
	}
}

// validateToken valida un token JWT
func validateToken(tokenString string) (*jwt.Token, error) {
	// Obtener la clave secreta desde las variables de entorno
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET no está configurado")
	}

	// Parsear y validar el token
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea el correcto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

// GenerateToken genera un nuevo token JWT para un usuario
func GenerateToken(userID int64) (string, error) {
	// Obtener la clave secreta desde las variables de entorno
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET no está configurado")
	}

	// Crear los claims del token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(jwt.TimeFunc().Add(jwt.DefaultLeeway)),
	}

	// Crear el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	return token.SignedString([]byte(secret))
}