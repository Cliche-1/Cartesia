package handlers

import (
	"Gin/views/layouts"
	"Gin/views/pages"
	"github.com/gin-gonic/gin"
)

// PageHandler maneja las rutas de páginas
type PageHandler struct{}

// NewPageHandler crea una nueva instancia de PageHandler
func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

// Home renderiza la página principal
func (h *PageHandler) Home(c *gin.Context) {
	title := "Inicio"
	content := pages.Home()
	component := layouts.Base(title, content)
	component.Render(c.Request.Context(), c.Writer)
}