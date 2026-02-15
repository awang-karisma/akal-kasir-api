package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func HandleResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func HandleError(w http.ResponseWriter, statusCode int, message string) {
	HandleResponse(w, statusCode, map[string]string{
		"status":  "error",
		"type":    http.StatusText(statusCode),
		"code":    strconv.Itoa(statusCode),
		"message": message,
	})
}

func IsDateValid(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
