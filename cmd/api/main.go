package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Gin/internal/database"
	"Gin/internal/handlers"
	"Gin/internal/middleware"
	"Gin/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No se pudo cargar el archivo .env: %v", err)
	}
}

func setupRouter(db *database.DB) *gin.Engine {
	// Crear router de Gin
	r := gin.Default()

	// Configurar middleware global
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Servir archivos estáticos
	r.Static("/static", "./static")

	// Inicializar handlers
	pageHandler := handlers.NewPageHandler(db)
	roadmapHandler := handlers.NewRoadmapHandler()
	nodeHandler := handlers.NewNodeHandler(db.GetDB())
	connectionHandler := handlers.NewConnectionHandler(db.GetDB())
	resourceHandler := handlers.NewResourceHandler(db.GetDB())

	// Inicializar servicios
	jwtService := services.NewJWTService()
	googleAuthService := services.NewGoogleAuthService()
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	ownerMiddleware := middleware.RequireRoadmapOwner(db.GetDB())

	// Rutas de páginas
	r.GET("/", pageHandler.Home)
	r.GET("/login", pageHandler.Login)
	r.GET("/register", pageHandler.Register)
	r.GET("/explore", pageHandler.Explore)

	// Rutas de autenticación
	authHandler := handlers.NewAuthHandler(db, jwtService, googleAuthService)
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", authMiddleware.RequireAuth(), authHandler.GetMe)
		auth.GET("/google/login", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
	}

	// Rutas de roadmaps
	roadmaps := r.Group("/roadmaps")
	{
		roadmaps.GET("/", roadmapHandler.ListRoadmaps)
		roadmaps.GET("/explore", roadmapHandler.ListRoadmaps)
		roadmaps.POST("/create", roadmapHandler.CreateRoadmap)
		
		roadmap := roadmaps.Group("/:id")
		{
			roadmap.GET("", roadmapHandler.ViewRoadmap)
			roadmap.PUT("", roadmapHandler.UpdateRoadmap)
			roadmap.DELETE("", roadmapHandler.DeleteRoadmap)
			roadmap.POST("/fork", roadmapHandler.ForkRoadmap)
			roadmap.GET("/reviews", roadmapHandler.GetRoadmapReviews)
			roadmap.POST("/reviews", roadmapHandler.AddReview)
			roadmap.GET("/editor", authMiddleware.RequireAuth(), ownerMiddleware, pageHandler.RoadmapEditor)
			
			// Rutas de nodos
			nodes := roadmap.Group("/nodes")
			{
				nodes.GET("/node/:node_id/resources", roadmapHandler.GetNodeResources)
				nodes.POST("/node/:node_id/complete", roadmapHandler.CompleteNode)
				nodes.POST("/node/:node_id/progress", roadmapHandler.UpdateProgress)
			}
		}
	}

	// Rutas de la API
	api := r.Group("/api")
	{
		// Rutas públicas
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"message": "El servidor está funcionando correctamente",
			})
		})

		// Rutas del editor de roadmaps (protegidas)
		apiRoadmaps := api.Group("/roadmaps/:id", authMiddleware.RequireAuth(), ownerMiddleware)
		{
			// Rutas de nodos
			apiRoadmaps.POST("/nodes", nodeHandler.CreateNode)
			apiRoadmaps.PUT("/nodes/:node_id", nodeHandler.UpdateNode)
			apiRoadmaps.DELETE("/nodes/:node_id", nodeHandler.DeleteNode)
			apiRoadmaps.PUT("/nodes/positions", nodeHandler.UpdateNodePositions)

			// Rutas de conexiones
			apiRoadmaps.POST("/connections", connectionHandler.CreateConnection)
			apiRoadmaps.DELETE("/connections/:conn_id", connectionHandler.DeleteConnection)

			// Rutas de recursos
			apiRoadmaps.POST("/nodes/:node_id/resources", resourceHandler.AddNodeResource)
		}
	}

	return r
}

func main() {
	// Configurar modo de Gin basado en variable de entorno
	if os.Getenv("ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar base de datos
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	// Obtener puerto del servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	// Inicializar router
	router := setupRouter(db)

	// Iniciar servidor
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor iniciando en http://localhost%s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}