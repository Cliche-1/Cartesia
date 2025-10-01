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
	ID          int64     `json:"id"`
	RoadmapID   int64     `json:"roadmap_id"`
	FromNodeID  int64     `json:"from_node_id"`
	ToNodeID    int64     `json:"to_node_id"`
	Label       string    `json:"label,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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