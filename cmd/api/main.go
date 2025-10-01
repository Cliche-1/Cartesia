package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Gin/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No se pudo cargar el archivo .env: %v", err)
	}
}

func setupRouter() *gin.Engine {
	// Crear router de Gin
	r := gin.Default()

	// Configurar middleware global
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Servir archivos estáticos
	r.Static("/static", "./static")

	// Inicializar handlers
	pageHandler := handlers.NewPageHandler()
	roadmapHandler := handlers.NewRoadmapHandler()

	// Rutas de páginas
	r.GET("/", pageHandler.Home)

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
			roadmap.POST("/reviews", roadmapHandler.AddReview)
			
			// Rutas de nodos
			nodes := roadmap.Group("/nodes")
			{
				nodes.POST("/:node_id/progress", roadmapHandler.UpdateProgress)
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

		// TODO: Agregar más rutas aquí
		// api.POST("/auth/login", handlers.Login)
		// api.POST("/auth/register", handlers.Register)

		// Rutas protegidas
		// authorized := api.Group("/")
		// authorized.Use(middleware.AuthRequired())
		// {
		//     authorized.GET("/profile", handlers.GetProfile)
		// }
	}

	return r
}

func main() {
	// Configurar modo de Gin basado en variable de entorno
	if os.Getenv("ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Obtener puerto del servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	// Inicializar router
	router := setupRouter()

	// Iniciar servidor
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor iniciando en http://localhost%s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}