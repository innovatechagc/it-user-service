package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"it-user-service/internal/config"
	"it-user-service/internal/models"
)

var DB *gorm.DB

// Connect establece la conexión con la base de datos
func Connect(cfg config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate ejecuta las migraciones automáticas para tablas de usuarios
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Solo migrar la tabla User que ya existe - no crear nuevas tablas
	// Las funcionalidades de perfiles, roles, etc. son placeholder por ahora
	err := DB.AutoMigrate(
		&models.User{}, // Tabla principal que ya existe en it-app_user
	)

	if err != nil {
		return fmt.Errorf("failed to migrate user database tables: %w", err)
	}

	log.Println("User service database migration completed successfully (using existing tables)")
	return nil
}

// Funciones de roles y permisos removidas - son placeholder por ahora
// Las funcionalidades de roles se implementarán más adelante cuando sea necesario

// Close cierra la conexión a la base de datos
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}