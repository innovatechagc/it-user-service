package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"it-user-service/internal/repositories"
	
)

// SetupRoutes configura todas las rutas del servicio
func SetupRoutes(userRepo repositories.UserRepositoryInterface, profileRepo repositories.ProfileRepositoryInterface, roleRepo repositories.RoleRepositoryInterface) *mux.Router {
	router := mux.NewRouter()

	// Crear handlers
	userHandler := NewUserHandler(userRepo)
	profileHandler := NewProfileHandler(profileRepo)
	roleHandler := NewRoleHandler(roleRepo)

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check routes
	api.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")

	// User routes
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/users/create", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	api.HandleFunc("/users/firebase/{firebase_id}", userHandler.GetUserByFirebaseID).Methods("GET")
	api.HandleFunc("/users/search", userHandler.SearchUsers).Methods("GET")

	// Profile routes
	api.HandleFunc("/users/{id}/profile", profileHandler.GetUserProfile).Methods("GET")
	api.HandleFunc("/users/{id}/profile", profileHandler.UpdateUserProfile).Methods("PUT")
	api.HandleFunc("/users/{id}/settings", profileHandler.GetUserSettings).Methods("GET")
	api.HandleFunc("/users/{id}/settings", profileHandler.UpdateUserSettings).Methods("PUT")
	api.HandleFunc("/users/{id}/stats", profileHandler.GetUserStats).Methods("GET")

	// Role routes
	api.HandleFunc("/roles", roleHandler.GetAllRoles).Methods("GET")
	api.HandleFunc("/roles/{id}", roleHandler.GetRoleByID).Methods("GET")
	api.HandleFunc("/roles", roleHandler.CreateRole).Methods("POST")
	api.HandleFunc("/roles/{id}", roleHandler.UpdateRole).Methods("PUT")
	api.HandleFunc("/roles/{id}", roleHandler.DeleteRole).Methods("DELETE")

	// User-Role assignment routes
	api.HandleFunc("/users/{user_id}/roles", roleHandler.AssignRoleToUser).Methods("POST")
	api.HandleFunc("/users/{user_id}/roles/{role_name}", roleHandler.RemoveRoleFromUser).Methods("DELETE")
	api.HandleFunc("/users/{user_id}/roles", roleHandler.GetUserRoles).Methods("GET")

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
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

	

	return router
}