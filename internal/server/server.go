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
	log := logger.GetLogger()

	// Conectar a la base de datos
	if err := database.Connect(cfg); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ejecutar migraciones
	if err := database.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Inicializar repositorios
	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	profileRepo := repositories.NewProfileRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	server := &Server{
		config:      cfg,
		router:      mux.NewRouter(),
		userRepo:    userRepo,
		profileRepo: profileRepo,
		roleRepo:    roleRepo,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// Crear handlers
	userHandler := handlers.NewUserHandler(s.userRepo)
	profileHandler := handlers.NewProfileHandler(s.profileRepo)
	roleHandler := handlers.NewRoleHandler(s.roleRepo)

	// Health check
	s.router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")

	// User routes
	s.router.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	s.router.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	s.router.HandleFunc("/users/create", userHandler.CreateUser).Methods("POST")
	s.router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	s.router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	s.router.HandleFunc("/users/firebase/{firebase_id}", userHandler.GetUserByFirebaseID).Methods("GET")
	s.router.HandleFunc("/users/search", userHandler.SearchUsers).Methods("GET")

	// Profile routes
	s.router.HandleFunc("/users/{id}/profile", profileHandler.GetUserProfile).Methods("GET")
	s.router.HandleFunc("/users/{id}/profile", profileHandler.UpdateUserProfile).Methods("PUT")
	s.router.HandleFunc("/users/{id}/settings", profileHandler.GetUserSettings).Methods("GET")
	s.router.HandleFunc("/users/{id}/settings", profileHandler.UpdateUserSettings).Methods("PUT")
	s.router.HandleFunc("/users/{id}/stats", profileHandler.GetUserStats).Methods("GET")

	// Role routes
	s.router.HandleFunc("/roles", roleHandler.GetAllRoles).Methods("GET")
	s.router.HandleFunc("/roles/{id}", roleHandler.GetRoleByID).Methods("GET")
	s.router.HandleFunc("/roles", roleHandler.CreateRole).Methods("POST")
	s.router.HandleFunc("/roles/{id}", roleHandler.UpdateRole).Methods("PUT")
	s.router.HandleFunc("/roles/{id}", roleHandler.DeleteRole).Methods("DELETE")

	// User-Role relationships
	s.router.HandleFunc("/users/{user_id}/roles", roleHandler.AssignRoleToUser).Methods("POST")
	s.router.HandleFunc("/users/{user_id}/roles/{role_id}", roleHandler.RemoveRoleFromUser).Methods("DELETE")
	s.router.HandleFunc("/users/{user_id}/roles", roleHandler.GetUserRoles).Methods("GET")
	s.router.HandleFunc("/users/{user_id}/permissions", roleHandler.GetUserPermissions).Methods("GET")

	// CORS middleware
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})

	// Logging middleware
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetLogger()
			log.WithFields(map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			}).Info("HTTP Request")
			next.ServeHTTP(w, r)
		})
	})
}

func (s *Server) Start() error {
	log := logger.GetLogger()
	
	addr := fmt.Sprintf(":%s", s.config.Port)
	log.WithField("address", addr).Info("Starting User Service server")
	
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) Close() error {
	return database.Close()
}