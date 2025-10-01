package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"Gin/internal/models"
	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	db *sql.DB
}

func NewConnectionHandler(db *sql.DB) *ConnectionHandler {
	return &ConnectionHandler{db: db}
}

// CreateConnection crea una nueva conexión entre nodos
func (h *ConnectionHandler) CreateConnection(c *gin.Context) {
	type createConnectionRequest struct {
		FromNodeID     int64                 `json:"from_node_id" binding:"required"`
		ToNodeID       int64                 `json:"to_node_id" binding:"required"`
		ConnectionType models.ConnectionType `json:"connection_type" binding:"required"`
	}

	var req createConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	roadmapID := c.GetInt64("roadmap_id")

	// Verificar que ambos nodos existen y pertenecen al roadmap
	var count int
	err := h.db.QueryRow(`
		SELECT COUNT(*)
		FROM nodes
		WHERE roadmap_id = $1 AND id IN ($2, $3)`,
		roadmapID, req.FromNodeID, req.ToNodeID,
	).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar nodos"})
		return
	}

	if count != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uno o ambos nodos no existen en este roadmap"})
		return
	}

	// Verificar que no existe ya una conexión entre estos nodos
	err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM connections
		WHERE roadmap_id = $1 AND from_node_id = $2 AND to_node_id = $3`,
		roadmapID, req.FromNodeID, req.ToNodeID,
	).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar conexión existente"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Ya existe una conexión entre estos nodos"})
		return
	}

	now := time.Now()

	// Crear la conexión
	var connectionID int64
	err = h.db.QueryRow(`
		INSERT INTO connections (roadmap_id, from_node_id, to_node_id, connection_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		roadmapID, req.FromNodeID, req.ToNodeID, req.ConnectionType, now, now,
	).Scan(&connectionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear conexión"})
		return
	}

	// Devolver la conexión creada
	connection := models.Connection{
		ID:             connectionID,
		RoadmapID:      roadmapID,
		FromNodeID:     req.FromNodeID,
		ToNodeID:       req.ToNodeID,
		ConnectionType: req.ConnectionType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	c.JSON(http.StatusCreated, connection)
}

// DeleteConnection elimina una conexión existente
func (h *ConnectionHandler) DeleteConnection(c *gin.Context) {
	connectionID, err := strconv.ParseInt(c.Param("conn_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de conexión inválido"})
		return
	}

	roadmapID := c.GetInt64("roadmap_id")

	// Eliminar la conexión
	result, err := h.db.Exec(`
		DELETE FROM connections
		WHERE id = $1 AND roadmap_id = $2`,
		connectionID, roadmapID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar conexión"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar eliminación"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conexión no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conexión eliminada correctamente"})
}