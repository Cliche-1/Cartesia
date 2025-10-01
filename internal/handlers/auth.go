package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"

	"Gin/internal/database"
	"Gin/internal/models"
	"Gin/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db              *database.DB
	jwtService      *services.JWTService
	googleAuthService *services.GoogleAuthService
}

func NewAuthHandler(db *database.DB, jwtService *services.JWTService, googleAuthService *services.GoogleAuthService) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtService:      jwtService,
		googleAuthService: googleAuthService,
	}
}

// generateState genera un estado aleatorio para OAuth2
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Register maneja el registro de nuevos usuarios
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro inválidos"})
		return
	}

	// Crear nuevo usuario
	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
	}

	// Validar usuario
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash de la contraseña
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insertar usuario en la base de datos
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := h.db.GetDB().QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint" {
			c.JSON(http.StatusConflict, gin.H{"error": "El usuario o email ya existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}

	// Generar token JWT
	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user,
	})
}

// Login maneja la autenticación de usuarios
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de login inválidos"})
		return
	}

	// Buscar usuario por email
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := h.db.GetDB().QueryRow(query, input.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	// Verificar contraseña
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Generar token JWT
	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// GoogleLogin inicia el flujo de autenticación con Google
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar estado"})
		return
	}

	// Guardar el estado en una cookie segura
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	// Redirigir a la URL de autorización de Google
	url := h.googleAuthService.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback maneja la respuesta de Google OAuth2
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Verificar el estado
	state, _ := c.Cookie("oauth_state")
	if state != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estado inválido"})
		return
	}

	// Obtener el token
	code := c.Query("code")
	token, err := h.googleAuthService.Exchange(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener token"})
		return
	}

	// Obtener información del usuario
	googleUser, err := h.googleAuthService.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	// Buscar o crear usuario
	var user models.User
	query := `
		SELECT id, username, email, avatar_url, created_at, updated_at 
		FROM users 
		WHERE google_id = $1 OR email = $2`

	err = h.db.GetDB().QueryRow(query, googleUser.ID, googleUser.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Crear nuevo usuario
		username := googleUser.GivenName + googleUser.FamilyName
		if username == "" {
			username = googleUser.Email[:strings.Index(googleUser.Email, "@")]
		}

		query = `
			INSERT INTO users (username, email, google_id, avatar_url)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`

		err = h.db.GetDB().QueryRow(
			query,
			username,
			googleUser.Email,
			googleUser.ID,
			googleUser.Picture,
		).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
			return
		}

		user.Username = username
		user.Email = googleUser.Email
		user.AvatarURL = googleUser.Picture
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	// Generar token JWT
	jwtToken, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	// Redirigir al frontend con el token
	c.Redirect(http.StatusTemporaryRedirect, "/auth/callback?token="+jwtToken)
}

// GetMe retorna la información del usuario actual
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Obtener user_id del contexto (establecido por el middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Buscar usuario en la base de datos
	var user models.User
	query := `SELECT id, username, email, avatar_url, bio, created_at, updated_at FROM users WHERE id = $1`
	err := h.db.GetDB().QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.AvatarURL,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}