package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	// Cargar variables desde .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se encontró el archivo .env, usando variables de entorno del sistema")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || dbname == "" {
		log.Fatal("Faltan variables de entorno para la conexión a la base de datos")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error conectando a la BD:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la BD:", err)
	}

	fmt.Println("Conectado a PostgreSQL exitosamente")
}
