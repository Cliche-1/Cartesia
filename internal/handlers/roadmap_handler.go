package handlers

import (
	"Gin/views/layouts"
	"Gin/views/pages"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoadmapHandler maneja las rutas relacionadas con roadmaps
type RoadmapHandler struct {
	// TODO: Agregar servicios necesarios
}

// NewRoadmapHandler crea una nueva instancia de RoadmapHandler
func NewRoadmapHandler() *RoadmapHandler {
	return &RoadmapHandler{}
}

// ListRoadmaps muestra la página principal con los roadmaps destacados
func (h *RoadmapHandler) ListRoadmaps(c *gin.Context) {
	// TODO: Obtener roadmaps destacados de la base de datos
	title := "Roadmaps Destacados"
	content := pages.Home() // Por ahora usamos la misma página de inicio
	component := layouts.Base(title, content)
	component.Render(c.Request.Context(), c.Writer)
}

// ViewRoadmap muestra un roadmap específico
func (h *RoadmapHandler) ViewRoadmap(c *gin.Context) {
	id := c.Param("id")
	// TODO: Obtener roadmap de la base de datos
	c.JSON(http.StatusOK, gin.H{
		"message": "Vista de roadmap implementada próximamente",
		"id":      id,
	})
}

// CreateRoadmap maneja la creación de un nuevo roadmap
func (h *RoadmapHandler) CreateRoadmap(c *gin.Context) {
	// TODO: Implementar creación de roadmap
	c.JSON(http.StatusOK, gin.H{
		"message": "Creación de roadmap implementada próximamente",
	})
}

// UpdateRoadmap maneja la actualización de un roadmap
func (h *RoadmapHandler) UpdateRoadmap(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implementar actualización de roadmap
	c.JSON(http.StatusOK, gin.H{
		"message": "Actualización de roadmap implementada próximamente",
		"id":      id,
	})
}

// DeleteRoadmap maneja la eliminación de un roadmap
func (h *RoadmapHandler) DeleteRoadmap(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implementar eliminación de roadmap
	c.JSON(http.StatusOK, gin.H{
		"message": "Eliminación de roadmap implementada próximamente",
		"id":      id,
	})
}

// ForkRoadmap crea una copia de un roadmap existente
func (h *RoadmapHandler) ForkRoadmap(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implementar fork de roadmap
	c.JSON(http.StatusOK, gin.H{
		"message": "Fork de roadmap implementado próximamente",
		"id":      id,
	})
}

// AddReview añade una reseña a un roadmap
func (h *RoadmapHandler) AddReview(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implementar añadir reseña
	c.JSON(http.StatusOK, gin.H{
		"message": "Añadir reseña implementado próximamente",
		"id":      id,
	})
}

// UpdateProgress actualiza el progreso en un nodo del roadmap
func (h *RoadmapHandler) UpdateProgress(c *gin.Context) {
	roadmapID := c.Param("roadmap_id")
	nodeID := c.Param("node_id")
	// TODO: Implementar actualización de progreso
	c.JSON(http.StatusOK, gin.H{
		"message":    "Actualización de progreso implementada próximamente",
		"roadmap_id": roadmapID,
		"node_id":    nodeID,
	})
}