package main

import (
	"cat-manager/backend/internal/handler"
	"cat-manager/backend/internal/middleware"
	"cat-manager/backend/internal/repository"
	"cat-manager/backend/internal/service"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

//go:embed all:static
var frontendFS embed.FS

func main() {
	host := flag.String("host", "", "listen IP (default: all interfaces)")
	port := flag.Int("port", 0, "listen port (default: 8080)")
	flag.Parse()

	addr := resolveAddr(*host, *port)

	// Ensure uploads directory exists
	os.MkdirAll("uploads/avatars", 0755)
	os.MkdirAll("uploads/photos", 0755)
	os.MkdirAll("uploads/reports", 0755)

	// Database
	db, err := initDB("cat-manager.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Repositories
	petRepo := repository.NewPetRepo(db)
	weightRepo := repository.NewWeightRepo(db)
	healthRepo := repository.NewHealthRepo(db)
	photoRepo := repository.NewPhotoRepo(db)

	// Services
	weightSvc := service.NewWeightService(weightRepo)
	petSvc := service.NewPetService(petRepo, weightSvc)
	healthSvc := service.NewHealthService(healthRepo)
	photoSvc := service.NewPhotoService(photoRepo)
	dashboardSvc := service.NewDashboardService(petRepo, weightRepo, healthRepo)

	// Handlers
	petHandler := handler.NewPetHandler(petSvc)
	weightHandler := handler.NewWeightHandler(weightSvc)
	healthHandler := handler.NewHealthHandler(healthSvc)
	photoHandler := handler.NewPhotoHandler(photoSvc, petSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/dashboard", dashboardHandler.Get)

		r.Route("/pets", func(r chi.Router) {
			r.Get("/", petHandler.List)
			r.Post("/", petHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", petHandler.Get)
				r.Put("/", petHandler.Update)
				r.Delete("/", petHandler.Delete)
				r.Post("/avatar", petHandler.UploadAvatar)

				// Weight records
				r.Get("/weights", weightHandler.List)
				r.Post("/weights", weightHandler.Create)
				r.Delete("/weights/{rid}", weightHandler.Delete)

				// Health records
				r.Get("/health", healthHandler.List)
				r.Post("/health", healthHandler.Create)
				r.Put("/health/{rid}", healthHandler.Update)
				r.Delete("/health/{rid}", healthHandler.Delete)

				// Photos
				r.Get("/photos", photoHandler.List)
				r.Post("/photos", photoHandler.Upload)
				r.Delete("/photos/{pid}", photoHandler.Delete)
			})
		})
	})

	// Static file serving for uploads
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Frontend static files
	staticFS, err := fs.Sub(frontendFS, "static")
	if err == nil {
		fileServer := http.FileServer(http.FS(staticFS))
		r.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path
			if path == "/" {
				path = "/index.html"
			}
			if _, err := fs.Stat(staticFS, path[1:]); err != nil {
				req.URL.Path = "/"
			}
			fileServer.ServeHTTP(w, req)
		})
	}

	fmt.Printf("🐱 Cat Manager starting on http://%s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func resolveAddr(flagHost string, flagPort int) string {
	host := flagHost
	if host == "" {
		host = os.Getenv("HOST")
	}

	port := flagPort
	if port == 0 {
		if envPort := os.Getenv("PORT"); envPort != "" {
			fmt.Sscanf(envPort, "%d", &port)
		}
	}
	if port == 0 {
		port = 8080
	}

	if host == "" {
		return fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func initDB(dsn string) (*repository.DB, error) {
	db, err := repository.NewDB(dsn)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(migrationSQL); err != nil {
		return nil, fmt.Errorf("run migration: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}
