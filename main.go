package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID     int    `json:"id"`
	Nama   string `json:"nama"`
	Harga  int    `json:"harga"`
	Jumlah int    `json:"jumlah"`
}

var semuaProduk = []Produk{
	{ID: 1, Nama: "Apple iPhone 17 ProMax", Harga: 27000000, Jumlah: 10},
	{ID: 2, Nama: "Nokia 3310", Harga: 100000, Jumlah: 5},
}

func getProdukByID(w http.ResponseWriter, id int) {
	for _, produk := range semuaProduk {
		if produk.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}
	}
	// produk tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Product not found",
	})
}

func updateProdukByID(w http.ResponseWriter, r *http.Request, id int) {
	var produkBaru Produk
	err1 := json.NewDecoder(r.Body).Decode(&produkBaru)
	if err1 != nil {
		log.Println(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}
	for i, produk := range semuaProduk {
		if produk.ID == id {
			produk.Nama = produkBaru.Nama
			produk.Harga = produkBaru.Harga
			produk.Jumlah = produkBaru.Jumlah
			semuaProduk[i] = produk
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Product not found",
	})
}
func deleteProdukByID(w http.ResponseWriter, id int) {
	for i, produk := range semuaProduk {
		if produk.ID == id {
			semuaProduk = append(semuaProduk[:i], semuaProduk[i+1:]...)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}
	}
	// produk tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": "Product not found",
	})
}
func main() {
	log.Println("Starting server...")
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Kasir API is running",
		})
	})

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(semuaProduk)
		} else if r.Method == "POST" {
			var produk Produk
			err := json.NewDecoder(r.Body).Decode(&produk)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"status":  "error",
					"message": "Invalid request body",
				})
				return
			}
			produk.ID = len(semuaProduk) + 1
			semuaProduk = append(semuaProduk, produk)
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		}
	})

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Invalid request path",
			})
			return
		}
		switch r.Method {
		case "GET":
			getProdukByID(w, id)
			return
		case "PUT":
			updateProdukByID(w, r, id)
			return
		case "DELETE":
			deleteProdukByID(w, id)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Method not allowed",
		})
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
