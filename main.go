package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type URL struct {
	OriginalURL string
	ShortURL    string
}

var urlStore = make(map[string]string) // In-memory URL store
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", RedirectHandler)
	fmt.Println("Starting server http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Could not start server", err)
		return
	}
}

func GenerateShortUrl() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request", http.StatusMethodNotAllowed)
		return
	}
	// get url from original body
	origianlURL := r.FormValue("url")
	if origianlURL == "" {
		http.Error(w, "URL required", http.StatusBadRequest)
		return
	}

	shortURL := GenerateShortUrl()
	urlStore[shortURL] = origianlURL

	fmt.Fprintf(w, "Shortened URL: http://localhost:8080/%s\n", shortURL)
}

// handler that takes short URL and redirect to the original

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] // get url from path
	originalURL, exists := urlStore[shortURL]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusBadRequest)
		return
	}

	// redirect
	http.Redirect(w, r, originalURL, http.StatusFound)
}
