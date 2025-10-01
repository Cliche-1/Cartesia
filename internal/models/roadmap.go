package models

import (
	"time"
)

// Roadmap representa un roadmap creado por un usuario
type Roadmap struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	AuthorID    int64     `json:"author_id"`
	IsPublic    bool      `json:"is_public"`
	ForkedFrom  *int64    `json:"forked_from,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Node representa un nodo en el roadmap
type Node struct {
	ID          int64     `json:"id"`
	RoadmapID   int64     `json:"roadmap_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        NodeType  `json:"type"`
	Position    Position  `json:"position"`
	Status      string    `json:"status"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NodeType representa el tipo de nodo en el roadmap
type NodeType string

const (
	NodeTypeTopic     NodeType = "topic"
	NodeTypeResource  NodeType = "resource"
	NodeTypeChallenge NodeType = "challenge"
	NodeTypeMilestone NodeType = "milestone"
)

// Position representa la posición de un nodo en el canvas
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Connection representa una conexión entre dos nodos
type Connection struct {
	ID             int64          `json:"id"`
	RoadmapID      int64          `json:"roadmap_id"`
	FromNodeID     int64          `json:"from_node_id"`
	ToNodeID       int64          `json:"to_node_id"`
	Label          string         `json:"label,omitempty"`
	ConnectionType ConnectionType `json:"connection_type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// ConnectionType representa el tipo de conexión entre nodos
type ConnectionType string

const (
	ConnectionTypeDefault ConnectionType = "default"
	ConnectionTypeStrong  ConnectionType = "strong"
	ConnectionTypeWeak    ConnectionType = "weak"
	ConnectionTypeDashed  ConnectionType = "dashed"
)

// Resource representa un recurso asociado a un nodo
type Resource struct {
	ID          int64     `json:"id"`
	NodeID      int64     `json:"node_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"` // url, video, document, etc.
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Review representa una reseña de un roadmap
type Review struct {
	ID        int64     `json:"id"`
	RoadmapID int64     `json:"roadmap_id"`
	UserID    int64     `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Progress representa el progreso de un usuario en un roadmap
type Progress struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	NodeID    int64     `json:"node_id"`
	Status    string    `json:"status"` // not_started, in_progress, completed
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Props para las vistas
type RoadmapDetailProps struct {
	ID          string
	Title       string
	Description string
	Author      struct {
		ID        string
		Name      string
		AvatarURL string
	}
	Stats struct {
		Views     int
		Forks     int
		Favorites int
	}
	Nodes     []RoadmapNodeProps
	Resources []ResourceProps
	Reviews   []ReviewProps
}

type RoadmapNodeProps struct {
	ID          string
	Title       string
	Description string
	Type        string
	PositionX   float64
	PositionY   float64
	Status      string
	Connections []struct {
		TargetID string
		Type     string
	}
}

type ResourceProps struct {
	ID          string
	Title       string
	Type        string // "link", "video", "document"
	URL         string
	Description string
}

type ReviewProps struct {
	ID        string
	UserName  string
	AvatarURL string
	Rating    int
	Comment   string
	CreatedAt string
}