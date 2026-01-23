package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Produk struct {
	ID     string `json:"id"`
	Nama   string `json:"nama"`
	Harga  int    `json:"harga"`
	Jumlah int    `json:"jumlah"`
}

type Kategori struct {
	ID        string `json:"id"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
}

var semuaProduk = []Produk{
	{ID: "b83628e5-6bad-4d7d-868e-32a9fe8ad84f", Nama: "Apple iPhone 17 ProMax", Harga: 27000000, Jumlah: 10},
	{ID: "3e1ca5e5-a3cb-4705-9311-456f1f263d96", Nama: "Nokia 3310", Harga: 100000, Jumlah: 5},
}

var semuaKategori = []Kategori{
	{ID: "1518d979-1590-475d-84b3-54c9423f1686", Nama: "Mobile", Deskripsi: "Produk mobile"},
	{ID: "8bdff33e-bf7f-451d-a244-e2311761cb26", Nama: "Laptop", Deskripsi: "Produk laptop"},
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

func handleResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	handleResponse(w, statusCode, map[string]string{
		"status":  "error",
		"type":    http.StatusText(statusCode),
		"code":    strconv.Itoa(statusCode),
		"message": message,
	})
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	handleResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Kasir API is running",
	})
}

func getProduk(w http.ResponseWriter, r *http.Request) {
	handleResponse(w, http.StatusOK, semuaProduk)
}

func createProduk(w http.ResponseWriter, r *http.Request) {
	var produk Produk
	err := json.NewDecoder(r.Body).Decode(&produk)
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	produk.ID = id.String()
	semuaProduk = append(semuaProduk, produk)
	handleResponse(w, http.StatusCreated, produk)
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid id")
		return
	}
	for _, produk := range semuaProduk {
		if produk.ID == id.String() {
			handleResponse(w, http.StatusOK, produk)
			return
		}
	}
	handleError(w, http.StatusNotFound, "Product not found")
}

func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}
	var produkBaru Produk
	err = json.NewDecoder(r.Body).Decode(&produkBaru)
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	for i, produk := range semuaProduk {
		if produk.ID == id.String() {
			produk.Nama = produkBaru.Nama
			produk.Harga = produkBaru.Harga
			produk.Jumlah = produkBaru.Jumlah
			semuaProduk[i] = produk
			handleResponse(w, http.StatusOK, produk)
			return
		}
	}
	handleError(w, http.StatusNotFound, "Product not found")
}

func deleteProdukByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}
	log.Println(id)
	for i, produk := range semuaProduk {
		if produk.ID == id.String() {
			semuaProduk = append(semuaProduk[:i], semuaProduk[i+1:]...)
			handleResponse(w, http.StatusOK, produk)
			return
		}
	}
	handleError(w, http.StatusNotFound, "Product not found")
}

func getKategori(w http.ResponseWriter, r *http.Request) {
	handleResponse(w, http.StatusOK, semuaKategori)
}

func getKategoriByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}
	for _, kategori := range semuaKategori {
		if kategori.ID == id.String() {
			handleResponse(w, http.StatusOK, kategori)
			return
		}
	}
	handleError(w, http.StatusNotFound, "Kategori not found")
}

func createKategori(w http.ResponseWriter, r *http.Request) {
	var kategori Kategori
	err := json.NewDecoder(r.Body).Decode(&kategori)
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	kategori.ID = id.String()
	semuaKategori = append(semuaKategori, kategori)
	handleResponse(w, http.StatusCreated, kategori)
}

func updateKategoriByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}
	var kategoriBaru Kategori
	err = json.NewDecoder(r.Body).Decode(&kategoriBaru)
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	for i, kategori := range semuaKategori {
		if kategori.ID == id.String() {
			kategori.Nama = kategoriBaru.Nama
			kategori.Deskripsi = kategoriBaru.Deskripsi
			semuaKategori[i] = kategori
			handleResponse(w, http.StatusOK, kategori)
			return
		}
	}
	handleError(w, http.StatusNotFound, "Kategori not found")
}

func deleteKategoriByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusBadRequest, "Invalid uuid")
		return
	}
	log.Println(id)
	for i, kategori := range semuaKategori {
		if kategori.ID == id.String() {
			semuaKategori = append(semuaKategori[:i], semuaKategori[i+1:]...)
			handleResponse(w, http.StatusOK, kategori)
			return
		}
	}
	handleError(w, http.StatusNotFound, "Kategori not found")
}

func main() {
	log.Println("Starting server...")
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/healthz", getHealthz).Methods("GET")

	r.HandleFunc("/api/produk", getProduk).Methods("GET")
	r.HandleFunc("/api/produk", createProduk).Methods("POST")
	r.HandleFunc("/api/produk/{id}", getProdukByID).Methods("GET")
	r.HandleFunc("/api/produk/{id}", updateProdukByID).Methods("PUT")
	r.HandleFunc("/api/produk/{id}", deleteProdukByID).Methods("DELETE")

	r.HandleFunc("/api/kategori", getKategori).Methods("GET")
	r.HandleFunc("/api/kategori", createKategori).Methods("POST")
	r.HandleFunc("/api/kategori/{id}", getKategoriByID).Methods("GET")
	r.HandleFunc("/api/kategori/{id}", updateKategoriByID).Methods("PUT")
	r.HandleFunc("/api/kategori/{id}", deleteKategoriByID).Methods("DELETE")

	log.Println("Listening to port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
