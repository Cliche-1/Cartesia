package handlers

import (
	"Gin/internal/models"
	"Gin/views/components"
	"Gin/views/layouts"
	"Gin/views/pages"
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type RoadmapHandler struct {
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
	roadmapID := c.Param("id")

	// TODO: Obtener datos del roadmap desde la base de datos
	roadmap := &models.RoadmapDetailProps{
		ID:          roadmapID,
		Title:       "Aprendiendo Go desde Cero",
		Description: "Una guía completa para aprender Go, desde los conceptos básicos hasta temas avanzados.",
		Author: struct {
			ID        string
			Name      string
			AvatarURL string
		}{
			ID:        "1",
			Name:      "John Doe",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=john",
		},
		Stats: struct {
			Views     int
			Forks     int
			Favorites int
		}{
			Views:     1234,
			Forks:     7,
			Favorites: 56,
		},
		Nodes: []models.RoadmapNodeProps{
			{
				ID:          "1",
				Title:       "Introducción a Go",
				Description: "Conceptos básicos del lenguaje Go",
				Type:        "topic",
				PositionX:   100,
				PositionY:   100,
				Status:      "completed",
				Connections: []struct {
					TargetID string
					Type     string
				}{
					{TargetID: "2", Type: "required"},
				},
			},
			{
				ID:          "2",
				Title:       "Estructuras de Control",
				Description: "If, switch, for y más",
				Type:        "topic",
				PositionX:   300,
				PositionY:   100,
				Status:      "in_progress",
				Connections: []struct {
					TargetID string
					Type     string
				}{},
			},
		},
	}

	content := pages.RoadmapDetail(*roadmap)
	component := layouts.Base(roadmap.Title, content)
	component.Render(c.Request.Context(), c.Writer)
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

// GetNodeResources obtiene los recursos de un nodo específico
func (h *RoadmapHandler) GetNodeResources(c *gin.Context) {
	// TODO: Obtener recursos desde la base de datos
	resources := []models.ResourceProps{
		{
			ID:          "1",
			Title:       "Documentación Oficial de Go",
			Type:        "link",
			URL:         "https://golang.org/doc",
			Description: "Documentación oficial y guías de Go",
		},
		{
			ID:          "2",
			Title:       "Tour of Go",
			Type:        "video",
			URL:         "https://tour.golang.org",
			Description: "Tutorial interactivo para aprender Go",
		},
	}

	component := components.ResourceCard(resources[0])
	component.Render(c.Request.Context(), c.Writer)
}

func (h *RoadmapHandler) GetRoadmapReviews(c *gin.Context) {
	// TODO: Obtener reviews desde la base de datos
	reviews := []models.ReviewProps{
		{
			ID:        "1",
			UserName:  "alice",
			Rating:    5,
			Comment:   "¡Excelente roadmap! Muy bien estructurado y fácil de seguir.",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=alice",
			CreatedAt: "hace 2 días",
		},
		{
			ID:        "2",
			UserName:  "bob",
			Rating:    4,
			Comment:   "Muy útil para principiantes. Podría tener más recursos avanzados.",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=bob",
			CreatedAt: "hace 5 días",
		},
	}

	component := components.ReviewCard(reviews[0])
	component.Render(c.Request.Context(), c.Writer)
}

func (h *RoadmapHandler) CompleteNode(c *gin.Context) {
	// TODO: Actualizar el estado del nodo en la base de datos
	success := true

	if success {
		c.JSON(http.StatusOK, gin.H{"status": "completed"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo marcar el nodo como completado"})
	}
}

// Helper functions para renderizar componentes
func renderResources(resources []models.ResourceProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, resource := range resources {
			err := components.ResourceCard(resource).Render(ctx, w)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func renderReviews(reviews []models.ReviewProps) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, review := range reviews {
			err := components.ReviewCard(review).Render(ctx, w)
			if err != nil {
				return err
			}
		}
		return nil
	})
}