package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter is a wrapper around http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs HTTP requests with method, path, status code, and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a custom ResponseWriter to capture status code
		rw := newResponseWriter(w)
		
		// Process the request
		next.ServeHTTP(rw, r)
		
		// Log the request
		duration := time.Since(start)
		statusEmoji := getStatusEmoji(rw.statusCode)
		methodEmoji := getMethodEmoji(r.Method)
		
		log.Printf("%s %s %s %s - %d %s (%v)",
			methodEmoji,
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.statusCode,
			statusEmoji,
			duration,
		)
	})
}

func getStatusEmoji(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "✅"
	case statusCode >= 300 && statusCode < 400:
		return "🔄"
	case statusCode >= 400 && statusCode < 500:
		return "⚠️"
	case statusCode >= 500:
		return "❌"
	default:
		return "❓"
	}
}

func getMethodEmoji(method string) string {
	switch method {
	case "GET":
		return "📄"
	case "POST":
		return "📝"
	case "PUT":
		return "✏️"
	case "DELETE":
		return "🗑️"
	case "PATCH":
		return "🔧"
	default:
		return "🔗"
	}
}