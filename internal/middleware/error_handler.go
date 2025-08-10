package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// ErrorContext holds error information with stack trace
type ErrorContext struct {
	Error      error
	StatusCode int
	StackTrace []string
	RequestURI string
	Method     string
}

// errorResponseWriter wraps http.ResponseWriter to capture errors
type errorResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	errorContext *ErrorContext
}

func newErrorResponseWriter(w http.ResponseWriter, requestURI, method string) *errorResponseWriter {
	return &errorResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		errorContext: &ErrorContext{
			RequestURI: requestURI,
			Method:     method,
		},
	}
}

func (erw *errorResponseWriter) WriteHeader(code int) {
	erw.statusCode = code
	erw.errorContext.StatusCode = code
	erw.ResponseWriter.WriteHeader(code)
}

func (erw *errorResponseWriter) Write(data []byte) (int, error) {
	// If this is an error response, capture the error messaget
	if erw.statusCode >= 400 && erw.errorContext.Error == nil {
		erw.errorContext.Error = errors.New(string(data))
		erw.errorContext.StackTrace = captureStackTrace(3) // Skip middleware frames
	}
	return erw.ResponseWriter.Write(data)
}

// ErrorHandlingMiddleware captures and logs detailed error information
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create error-aware response writer
		erw := newErrorResponseWriter(w, r.RequestURI, r.Method)

		// Defer error logging
		defer func() {
			if rec := recover(); rec != nil {
				// Handle panics
				erw.errorContext.Error = fmt.Errorf("panic: %v", rec)
				erw.errorContext.StatusCode = http.StatusInternalServerError
				erw.errorContext.StackTrace = captureStackTrace(0)
				logError(erw.errorContext)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Log errors for 4xx and 5xx status codes
			if erw.statusCode >= 400 {
				logError(erw.errorContext)
			}
		}()

		// Process request
		next.ServeHTTP(erw, r)
	})
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) []string {
	var stackTrace []string

	// Get stack trace
	for i := skip; i < skip+10; i++ { // Capture up to 10 frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Get function name
		fn := runtime.FuncForPC(pc)
		var funcName string
		if fn != nil {
			funcName = fn.Name()
		} else {
			funcName = "unknown"
		}

		// Format stack frame
		frame := fmt.Sprintf("%s:%d %s", shortenPath(file), line, shortenFuncName(funcName))
		stackTrace = append(stackTrace, frame)
	}

	return stackTrace
}

// logError logs detailed error information
func logError(ctx *ErrorContext) {
	if ctx == nil || ctx.Error == nil {
		return
	}

	statusEmoji := getErrorEmoji(ctx.StatusCode)

	log.Printf("\n"+
		"🚨 ERROR TRACE %s\n"+
		"├─ Request: %s %s\n"+
		"├─ Status: %d %s\n"+
		"├─ Error: %s\n"+
		"└─ Stack Trace:",
		statusEmoji, ctx.Method, ctx.RequestURI, ctx.StatusCode, statusEmoji, ctx.Error.Error())

	for i, frame := range ctx.StackTrace {
		if i == len(ctx.StackTrace)-1 {
			log.Printf("   └─ %s", frame)
		} else {
			log.Printf("   ├─ %s", frame)
		}
	}
	log.Println() // Empty line for readability
}

// getErrorEmoji returns emoji based on status code
func getErrorEmoji(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "💥"
	case statusCode >= 400:
		return "⚠️"
	default:
		return "❓"
	}
}

// shortenPath shortens file paths for better readability
func shortenPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return ".../" + strings.Join(parts[len(parts)-3:], "/")
	}
	return path
}

// shortenFuncName shortens function names for better readability
func shortenFuncName(funcName string) string {
	parts := strings.Split(funcName, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		// Remove package path, keep only package.function
		if dotIndex := strings.LastIndex(lastPart, "."); dotIndex != -1 {
			return lastPart
		}
		return lastPart
	}
	return funcName
}
