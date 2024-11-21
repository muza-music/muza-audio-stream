package main

import (
	"fmt"
	"gas/pkg/auth"
	"gas/pkg/ffmpeg"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Serve the HTML client at the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	// Handle the /audio endpoint
	http.HandleFunc("/audio", AudioHandler)

	fmt.Println("Starting server on :8443")
	err := http.ListenAndServeTLS(":8443", "certs/server/cert.pem", "certs/server/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS error:", err)
	}
}

// AudioHandler handles requests to stream audio files
func AudioHandler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight request
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		return
	}

	// Allow cross-origin requests (CORS)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Only allow GET requests
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract JWT token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Bearer token missing", http.StatusUnauthorized)
		return
	}

	// Validate JWT token
	claims, err := auth.ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Get audio parameters from query parameters
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename parameter is required", http.StatusBadRequest)
		return
	}

	audioOpts := ffmpeg.AudioOptions{
		Filename:   filename,
		Bitrate:    r.URL.Query().Get("bitrate"),
		SampleRate: r.URL.Query().Get("samplerate"),
		Channels:   r.URL.Query().Get("channels"),
		Codec:      r.URL.Query().Get("codec"),
		Quality:    r.URL.Query().Get("quality"),
	}

	// Determine Content-Type based on codec
	contentType := "application/octet-stream" // Default if codec is unknown
	switch audioOpts.Codec {
	case "mp3":
		contentType = "audio/mpeg"
	case "aac":
		contentType = "audio/aac"
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", contentType)

	// Log start streaming
	log.Printf("Start streaming audio file '%s' to user '%s' ('%s') with options %+v\n", filename, claims.Username, claims.Phone, audioOpts)

	err = ffmpeg.StreamAudio(audioOpts, w)
	if err != nil {
		// Log error streaming
		log.Printf("Error streaming audio file '%s' to user '%s' ('%s'): %+v\n", filename, claims.Username, claims.Phone, err)

		http.Error(w, "Error streaming audio", http.StatusInternalServerError)
		return
	}

	// Log end streaming
	log.Printf("End   streaming audio file '%s' to user '%s' ('%s')\n", filename, claims.Username, claims.Phone)
}
