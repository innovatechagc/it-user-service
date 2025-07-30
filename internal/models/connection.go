package models

import (
	"log"
	"gorm.io/gorm"
	"it-user-service/internal/database"
)

// ConnectDB establece la conexión con la base de datos PostgreSQL
func ConnectDB() {
	database.ConnectDB()
}

// MigrateDB ejecuta las migraciones automáticas
func MigrateDB() {
	db := database.GetDB()
	
	// Habilitar extensión UUID si no existe
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	
	err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&UserSettings{},
		&UserStats{},
		&Role{},
		&UserRole{},
	)
	
	if err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}
	
	log.Println("Migraciones ejecutadas exitosamente")
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return database.GetDB()
}