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

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)
	healthHandler := handlers.NewHealthHandler(db)

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/healthz", healthHandler.HandleHealth)

	r.HandleFunc("/api/products", productHandler.HandleProduct)
	r.HandleFunc("/api/products/{id}", productHandler.HandleProductByID)
	r.HandleFunc("/api/products/{id}/categories", productHandler.HandleProductCategories)

	r.HandleFunc("/api/categories", categoryHandler.HandleCategory)
	r.HandleFunc("/api/categories/{id}", categoryHandler.HandleCategoryByID)
	r.HandleFunc("/api/categories/{id}/products", categoryHandler.GetProductsByCategory)

	r.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)
	r.HandleFunc("/api/transactions", transactionHandler.GetTransactions)

	r.HandleFunc("/api/reports", reportHandler.HandleReport)
	r.HandleFunc("/api/reports/today", reportHandler.GetReportToday)

	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		internal.HandleError(w, http.StatusNotFound, "Not found")
	})

	log.Println("Listening to port " + config.Port)
	err = http.ListenAndServe(":"+config.Port, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
