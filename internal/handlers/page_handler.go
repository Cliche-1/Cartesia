package handlers

import (
	"Gin/internal/database"
	"Gin/views/components"
	"Gin/views/layouts"
	"Gin/views/pages"

	"github.com/gin-gonic/gin"
)

// PageHandler maneja las rutas de páginas
type PageHandler struct {
	db *database.DB
}

// NewPageHandler crea una nueva instancia de PageHandler
func NewPageHandler(db *database.DB) *PageHandler {
	return &PageHandler{
		db: db,
	}
}

// Home renderiza la página principal
func (h *PageHandler) Home(c *gin.Context) {
	component := layouts.Base("Inicio - Cartesia", pages.Home())
	component.Render(c.Request.Context(), c.Writer)
}

// Explore renderiza la página de exploración de roadmaps
func (h *PageHandler) Explore(c *gin.Context) {
	// Obtener roadmaps de la base de datos
	rows, err := h.db.GetDB().Query(`
		SELECT r.id, r.title, r.description, r.created_at,
			   u.id as author_id, u.username as author_name, u.avatar_url,
			   COALESCE(v.views_count, 0) as views_count,
			   COALESCE(rv.avg_rating, 0) as avg_rating,
			   COALESCE(rv.reviews_count, 0) as reviews_count
		FROM roadmaps r
		LEFT JOIN users u ON r.author_id = u.id
		LEFT JOIN (
			SELECT roadmap_id, COUNT(*) as views_count
			FROM roadmap_views
			GROUP BY roadmap_id
		) v ON r.id = v.roadmap_id
		LEFT JOIN (
			SELECT roadmap_id,
				   AVG(rating) as avg_rating,
				   COUNT(*) as reviews_count
			FROM reviews
			GROUP BY roadmap_id
		) rv ON r.id = rv.roadmap_id
		WHERE r.is_public = true
		ORDER BY r.created_at DESC
		LIMIT 12
	`)
	if err != nil {
		c.Status(500)
		return
	}
	defer rows.Close()

	var roadmaps []components.RoadmapCardProps
	for rows.Next() {
		var r components.RoadmapCardProps
		err := rows.Scan(
			&r.ID, &r.Title, &r.Description, &r.ImageURL,
			&r.Author.ID, &r.Author.Name, &r.Author.AvatarURL,
			&r.Stats.Views, &r.Stats.Rating, &r.Stats.Reviews,
		)
		if err != nil {
			continue
		}
		roadmaps = append(roadmaps, r)
	}

	// Obtener categorías y tags para los filtros
	categories, _ := h.getCategories()
	tags, _ := h.getTags()

	props := pages.ExplorePageProps{
		Roadmaps: roadmaps,
		Filters: components.FiltersSidebarProps{
			Categories: categories,
			Tags:       tags,
		},
	}

	component := layouts.Base("Explorar Roadmaps - Cartesia", pages.ExplorePage(props))
	component.Render(c.Request.Context(), c.Writer)
}

func (h *PageHandler) getCategories() ([]struct {
	ID    string
	Name  string
	Count int
}, error) {
	rows, err := h.db.GetDB().Query(`
		SELECT c.id, c.name, COUNT(r.id) as count
		FROM categories c
		LEFT JOIN roadmap_categories rc ON c.id = rc.category_id
		LEFT JOIN roadmaps r ON rc.roadmap_id = r.id AND r.is_public = true
		GROUP BY c.id, c.name
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []struct {
		ID    string
		Name  string
		Count int
	}
	for rows.Next() {
		var c struct {
			ID    string
			Name  string
			Count int
		}
		if err := rows.Scan(&c.ID, &c.Name, &c.Count); err != nil {
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (h *PageHandler) getTags() ([]struct {
	ID    string
	Name  string
	Count int
}, error) {
	rows, err := h.db.GetDB().Query(`
		SELECT t.id, t.name, COUNT(r.id) as count
		FROM tags t
		LEFT JOIN roadmap_tags rt ON t.id = rt.tag_id
		LEFT JOIN roadmaps r ON rt.roadmap_id = r.id AND r.is_public = true
		GROUP BY t.id, t.name
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []struct {
		ID    string
		Name  string
		Count int
	}
	for rows.Next() {
		var t struct {
			ID    string
			Name  string
			Count int
		}
		if err := rows.Scan(&t.ID, &t.Name, &t.Count); err != nil {
			continue
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (h *PageHandler) Login(c *gin.Context) {
	component := layouts.Base("Iniciar Sesión - Cartesia", pages.LoginForm())
	component.Render(c.Request.Context(), c.Writer)
}

func (h *PageHandler) Register(c *gin.Context) {
	component := layouts.Base("Registrarse - Cartesia", pages.RegisterForm())
	component.Render(c.Request.Context(), c.Writer)
}

func (h *PageHandler) RoadmapEditor(c *gin.Context) {
	roadmapID := c.Param("id")

	// Obtener datos del roadmap desde la base de datos
	row := h.db.GetDB().QueryRow(`
		SELECT r.id, r.title, r.description, r.created_at,
			   u.id as author_id, u.username as author_name, u.avatar_url
		FROM roadmaps r
		LEFT JOIN users u ON r.author_id = u.id
		WHERE r.id = $1
	`, roadmapID)

	var roadmap struct {
		ID          string
		Title       string
		Description string
		CreatedAt   string
		Author      struct {
			ID        string
			Name      string
			AvatarURL string
		}
	}

	err := row.Scan(
		&roadmap.ID, &roadmap.Title, &roadmap.Description, &roadmap.CreatedAt,
		&roadmap.Author.ID, &roadmap.Author.Name, &roadmap.Author.AvatarURL,
	)
	if err != nil {
		c.Status(404)
		return
	}

	component := layouts.Base(roadmap.Title+" - Editor", pages.RoadmapEditor(roadmap.Title+" - Editor", roadmap.ID))
	component.Render(c.Request.Context(), c.Writer)
}