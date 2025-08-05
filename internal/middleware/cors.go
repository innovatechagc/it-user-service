package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		environment := os.Getenv("ENVIRONMENT")
		
		// Configuración de CORS según el ambiente
		if environment == "production" {
			// En producción, solo permitir orígenes específicos
			allowedOrigins := []string{
				"https://innovatech.com",
				"https://www.innovatech.com",
				"https://app.innovatech.com",
			}
			
			originAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					originAllowed = true
					break
				}
			}
			
			if !originAllowed && origin != "" {
				// Origen no permitido en producción
				w.WriteHeader(http.StatusForbidden)
				return
			}
		} else {
			// En desarrollo/staging, permitir orígenes locales
			if origin != "" {
				if strings.HasPrefix(origin, "http://localhost:") || 
				   strings.HasPrefix(origin, "http://127.0.0.1:") ||
				   strings.HasPrefix(origin, "https://localhost:") || 
				   strings.HasPrefix(origin, "https://127.0.0.1:") {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					// Para otros orígenes en desarrollo, usar lista específica
					devAllowedOrigins := []string{
						"https://staging.innovatech.com",
						"https://dev.innovatech.com",
					}
					
					for _, allowed := range devAllowedOrigins {
						if origin == allowed {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			} else {
				// Si no hay Origin header en desarrollo, permitir
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}
		
		// Headers CORS estándar
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		// Headers adicionales para evitar problemas
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// Handle preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}