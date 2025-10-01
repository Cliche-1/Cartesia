package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"Gin/internal/models"
	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	db *sql.DB
}

func NewNodeHandler(db *sql.DB) *NodeHandler {
	return &NodeHandler{db: db}
}

// CreateNode crea un nuevo nodo en el roadmap
func (h *NodeHandler) CreateNode(c *gin.Context) {
	type createNodeRequest struct {
		Title       string       `json:"title" binding:"required"`
		Description string       `json:"description"`
		Type        models.NodeType `json:"node_type" binding:"required"`
		PositionX   float64     `json:"position_x" binding:"required"`
		PositionY   float64     `json:"position_y" binding:"required"`
		Color       string      `json:"color" binding:"required"`
	}

	var req createNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	roadmapID := c.GetInt64("roadmap_id")
	now := time.Now()

	// Insertar el nuevo nodo
	var nodeID int64
	err := h.db.QueryRow(`
		INSERT INTO nodes (roadmap_id, title, description, type, position_x, position_y, color, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		roadmapID, req.Title, req.Description, req.Type, req.PositionX, req.PositionY, req.Color, "not_started", now, now,
	).Scan(&nodeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el nodo"})
		return
	}

	// Devolver el nodo creado
	node := models.Node{
		ID:          nodeID,
		RoadmapID:   roadmapID,
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Position:    models.Position{X: req.PositionX, Y: req.PositionY},
		Status:      "not_started",
		Color:       req.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	c.JSON(http.StatusCreated, node)
}

// UpdateNode actualiza un nodo existente
func (h *NodeHandler) UpdateNode(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de nodo inválido"})
		return
	}

	type updateNodeRequest struct {
		Title       *string       `json:"title"`
		Description *string       `json:"description"`
		Type        *models.NodeType `json:"node_type"`
		PositionX   *float64     `json:"position_x"`
		PositionY   *float64     `json:"position_y"`
		Color       *string      `json:"color"`
	}

	var req updateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Obtener el nodo actual
	var node models.Node
	err = h.db.QueryRow(`
		SELECT id, roadmap_id, title, description, type, position_x, position_y, status, color, created_at, updated_at
		FROM nodes WHERE id = $1 AND roadmap_id = $2`,
		nodeID, c.GetInt64("roadmap_id"),
	).Scan(
		&node.ID, &node.RoadmapID, &node.Title, &node.Description, &node.Type,
		&node.Position.X, &node.Position.Y, &node.Status, &node.Color, &node.CreatedAt, &node.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nodo no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el nodo"})
		return
	}

	// Actualizar solo los campos proporcionados
	if req.Title != nil {
		node.Title = *req.Title
	}
	if req.Description != nil {
		node.Description = *req.Description
	}
	if req.Type != nil {
		node.Type = *req.Type
	}
	if req.PositionX != nil {
		node.Position.X = *req.PositionX
	}
	if req.PositionY != nil {
		node.Position.Y = *req.PositionY
	}
	if req.Color != nil {
		node.Color = *req.Color
	}
	node.UpdatedAt = time.Now()

	// Actualizar el nodo en la base de datos
	_, err = h.db.Exec(`
		UPDATE nodes
		SET title = $1, description = $2, type = $3, position_x = $4, position_y = $5, color = $6, updated_at = $7
		WHERE id = $8 AND roadmap_id = $9`,
		node.Title, node.Description, node.Type, node.Position.X, node.Position.Y, node.Color, node.UpdatedAt,
		node.ID, node.RoadmapID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el nodo"})
		return
	}

	c.JSON(http.StatusOK, node)
}

// DeleteNode elimina un nodo y sus conexiones
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	nodeID, err := strconv.ParseInt(c.Param("node_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de nodo inválido"})
		return
	}

	roadmapID := c.GetInt64("roadmap_id")

	// Iniciar transacción
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al iniciar transacción"})
		return
	}
	defer tx.Rollback()

	// Eliminar las conexiones relacionadas con el nodo
	_, err = tx.Exec(`
		DELETE FROM connections
		WHERE roadmap_id = $1 AND (from_node_id = $2 OR to_node_id = $2)`,
		roadmapID, nodeID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar conexiones"})
		return
	}

	// Eliminar los recursos del nodo
	_, err = tx.Exec(`
		DELETE FROM resources
		WHERE node_id = $1`,
		nodeID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar recursos"})
		return
	}

	// Eliminar el progreso relacionado con el nodo
	_, err = tx.Exec(`
		DELETE FROM progress
		WHERE node_id = $1`,
		nodeID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar progreso"})
		return
	}

	// Eliminar el nodo
	result, err := tx.Exec(`
		DELETE FROM nodes
		WHERE id = $1 AND roadmap_id = $2`,
		nodeID, roadmapID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar nodo"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar eliminación"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nodo no encontrado"})
		return
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar transacción"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nodo eliminado correctamente"})
}

// UpdateNodePositions actualiza las posiciones de múltiples nodos
func (h *NodeHandler) UpdateNodePositions(c *gin.Context) {
	type nodePosition struct {
		NodeID    int64   `json:"node_id" binding:"required"`
		PositionX float64 `json:"position_x" binding:"required"`
		PositionY float64 `json:"position_y" binding:"required"`
	}

	var positions []nodePosition
	if err := c.ShouldBindJSON(&positions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	roadmapID := c.GetInt64("roadmap_id")
	now := time.Now()

	// Iniciar transacción
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al iniciar transacción"})
		return
	}
	defer tx.Rollback()

	// Actualizar cada posición
	for _, pos := range positions {
		_, err := tx.Exec(`
			UPDATE nodes
			SET position_x = $1, position_y = $2, updated_at = $3
			WHERE id = $4 AND roadmap_id = $5`,
			pos.PositionX, pos.PositionY, now, pos.NodeID, roadmapID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar posiciones"})
			return
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al confirmar transacción"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Posiciones actualizadas correctamente"})
}