package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"it-user-service/internal/config"
	"it-user-service/internal/database"
	"it-user-service/internal/handlers"
	"it-user-service/internal/logger"
	"it-user-service/internal/repositories"
)

type Server struct {
	config      config.Config
	router      *mux.Router
	userRepo    repositories.UserRepositoryInterface
	profileRepo repositories.ProfileRepositoryInterface
	roleRepo    repositories.RoleRepositoryInterface
}

func NewServer(cfg config.Config) (*Server, error) {
	// Conectar a la base de datos
	database.ConnectDB()

	// Inicializar repositorios
	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	profileRepo := repositories.NewProfileRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	server := &Server{
		config:      cfg,
		userRepo:    userRepo,
		profileRepo: profileRepo,
		roleRepo:    roleRepo,
	}

	server.router = handlers.SetupRoutes(server.userRepo, server.profileRepo, server.roleRepo)
	return server, nil
}

func (s *Server) Start() error {
	log := logger.GetLogger()
	
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.WithField("address", addr).Info("Starting User Service server")
	
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) Close() error {
	sqlDB, err := database.GetDB().DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}