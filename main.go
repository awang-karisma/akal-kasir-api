package main

import (
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/internal"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		log.Printf("[REQUEST] %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("[RESPONSE] %s %s %s %d %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func main() {
	log.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s. Using environment variables.", err)
	}
	var config = Config{
		Port:   os.Getenv("PORT"),
		DBConn: os.Getenv("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal("Error closing database: ", err)
		}
	}()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	healthHandler := handlers.NewHealthHandler(db)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/healthz", healthHandler.HandleHealth).Methods(http.MethodGet)

	r.HandleFunc("/api/products", productHandler.HandleProduct).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/api/products/{id}", productHandler.HandleProductByID).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		internal.HandleError(w, http.StatusNotFound, "Not found")
	})

	log.Println("Listening to port " + config.Port)
	err = http.ListenAndServe(":"+config.Port, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
