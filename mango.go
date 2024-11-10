package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urlMap = make(map[string]string)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow only your frontend origin
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}


func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	
	var requestBody struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.URL == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}


	shortCode := generateShortCode()
	urlMap[shortCode] = requestBody.URL 


	response := map[string]string{"short_url": fmt.Sprintf("http://localhost:8080/r/%s", shortCode)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func redirectHandler(w http.ResponseWriter, r *http.Request) {
	
	shortCode := strings.TrimPrefix(r.URL.Path, "/r/")
	originalURL, exists := urlMap[shortCode]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}


	http.Redirect(w, r, originalURL, http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {

	http.HandleFunc("/", indexHandler)


	http.HandleFunc("/shorten", shortenURLHandler)


	http.HandleFunc("/r/", redirectHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
