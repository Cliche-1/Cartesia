package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() (*sql.DB, error) {
	// Obtener variables de entorno
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construir string de conexión
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Abrir conexión
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión con la base de datos: %v", err)
	}

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error al verificar conexión con la base de datos: %v", err)
	}

	log.Println("Conexión exitosa con la base de datos")
	return db, nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *sql.DB {
	return db
}

// CloseDB cierra la conexión con la base de datos
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}