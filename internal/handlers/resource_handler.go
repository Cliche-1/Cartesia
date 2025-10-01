package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"Gin/internal/models"
	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	db *sql.DB
}

func NewResourceHandler(db *sql.DB) *ResourceHandler {
	return &ResourceHandler{db: db}
}

// AddNodeResource añade un nuevo recurso a un nodo
func (h *ResourceHandler) AddNodeResource(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de nodo inválido"})
		return
	}

	type addResourceRequest struct {
		ResourceType string `json:"resource_type" binding:"required"`
		Title       string `json:"title" binding:"required"`
		URL         string `json:"url" binding:"required"`
		Description string `json:"description"`
	}

	var req addResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	roadmapID := c.GetInt64("roadmap_id")

	// Verificar que el nodo existe y pertenece al roadmap
	var exists bool
	err = h.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM nodes
			WHERE id = $1 AND roadmap_id = $2
		)`,
		nodeID, roadmapID,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar nodo"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nodo no encontrado"})
		return
	}

	now := time.Now()

	// Crear el recurso
	var resourceID int64
	err = h.db.QueryRow(`
		INSERT INTO resources (node_id, title, type, url, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		nodeID, req.Title, req.ResourceType, req.URL, req.Description, now, now,
	).Scan(&resourceID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear recurso"})
		return
	}

	// Devolver el recurso creado
	resource := models.Resource{
		ID:          resourceID,
		NodeID:      nodeID,
		Title:       req.Title,
		Type:        req.ResourceType,
		URL:         req.URL,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	c.JSON(http.StatusCreated, resource)
}