package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Produk struct {
	ID     string `json:"id"`
	Nama   string `json:"nama"`
	Harga  int    `json:"harga"`
	Jumlah int    `json:"jumlah"`
}

var semuaProduk = []Produk{
	{ID: "b83628e5-6bad-4d7d-868e-32a9fe8ad84f", Nama: "Apple iPhone 17 ProMax", Harga: 27000000, Jumlah: 10},
	{ID: "3e1ca5e5-a3cb-4705-9311-456f1f263d96", Nama: "Nokia 3310", Harga: 100000, Jumlah: 5},
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
	log.Println("GET produk")
	handleResponse(w, http.StatusOK, semuaProduk)
}

func createProduk(w http.ResponseWriter, r *http.Request) {
	log.Println("POST Produk")
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
	log.Println("GET produk by id")
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
	log.Println("PUT produk")
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
	log.Printf("DELETE produk by id %s", id.String())

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

func main() {
	log.Println("Starting server...")
	r := mux.NewRouter()
	r.HandleFunc("/healthz", getHealthz).Methods("GET")
	r.HandleFunc("/api/produk", getProduk).Methods("GET")
	r.HandleFunc("/api/produk", createProduk).Methods("POST")
	r.HandleFunc("/api/produk/{id}", getProdukByID).Methods("GET")
	r.HandleFunc("/api/produk/{id}", updateProdukByID).Methods("PUT")
	r.HandleFunc("/api/produk/{id}", deleteProdukByID).Methods("DELETE")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	})

	log.Println("Listening to port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
