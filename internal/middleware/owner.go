package middleware

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RequireRoadmapOwner verifica que el usuario autenticado es el propietario del roadmap
func RequireRoadmapOwner(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID del usuario autenticado del contexto
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			c.Abort()
			return
		}

		// Obtener el ID del roadmap de los parámetros
		roadmapID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de roadmap inválido"})
			c.Abort()
			return
		}

		// Verificar si el usuario es el propietario del roadmap
		var authorID int64
		err = db.QueryRow("SELECT author_id FROM roadmaps WHERE id = $1", roadmapID).Scan(&authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Roadmap no encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar propiedad del roadmap"})
			}
			c.Abort()
			return
		}

		if authorID != userID.(int64) {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para modificar este roadmap"})
			c.Abort()
			return
		}

		// Almacenar el ID del roadmap en el contexto para uso posterior
		c.Set("roadmap_id", roadmapID)
		c.Next()
	}
}