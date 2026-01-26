package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Struktur log entry
type APILogEntry struct {
	Timestamp    string            `json:"timestamp"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	Query        map[string]string `json:"query,omitempty"`
	Headers      map[string]string `json:"headers"`
	RequestBody  interface{}       `json:"request_body,omitempty"`
	StatusCode   int               `json:"status_code"`
	ResponseBody interface{}       `json:"response_body,omitempty"`
	Duration     string            `json:"duration"`
	DurationMs   int64             `json:"duration_ms"`
	ClientIP     string            `json:"client_ip"`
	UserAgent    string            `json:"user_agent"`
	UserID       interface{}       `json:"user_id,omitempty"`
	UserEmail    interface{}       `json:"user_email,omitempty"`
	Error        string            `json:"error,omitempty"`
}

// Custom response writer untuk capture response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Main middleware function
func APILoggerMiddleware() gin.HandlerFunc {
	// Buat folder logs jika belum ada
	logsDir := "./logs"
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
	}

	return func(c *gin.Context) {
		startTime := time.Now()

		// === 1. CAPTURE REQUEST BODY ===
		var requestBody interface{}
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				// Restore body untuk handler
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Parse JSON jika bisa
				var jsonBody interface{}
				if json.Unmarshal(bodyBytes, &jsonBody) == nil {
					requestBody = jsonBody
				} else {
					requestBody = string(bodyBytes)
				}
			}
		}

		// === 2. SETUP RESPONSE CAPTURE ===
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseWriter

		// === 3. PROCESS REQUEST ===
		c.Next()

		// === 4. CALCULATE DURATION ===
		duration := time.Since(startTime)

		// === 5. PARSE RESPONSE BODY ===
		var responseBody interface{}
		if responseWriter.body.Len() > 0 {
			var jsonResponse interface{}
			if json.Unmarshal(responseWriter.body.Bytes(), &jsonResponse) == nil {
				responseBody = jsonResponse
			} else {
				responseBody = responseWriter.body.String()
			}
		}

		// === 6. GET QUERY PARAMS ===
		queryParams := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				queryParams[key] = values[0]
			}
		}

		// === 7. GET HEADERS (IMPORTANT ONES) ===
		headers := map[string]string{
			"Content-Type":  c.GetHeader("Content-Type"),
			"Accept":        c.GetHeader("Accept"),
			"Authorization": maskAuthHeader(c.GetHeader("Authorization")),
		}

		// === 8. GET USER INFO (IF AUTHENTICATED) ===
		var userID, userEmail interface{}
		if id, exists := c.Get("user_id"); exists {
			userID = id
		}
		if email, exists := c.Get("user_email"); exists {
			userEmail = email
		}

		// === 9. GET ERROR IF ANY ===
		var errorMsg string
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// === 10. CREATE LOG ENTRY ===
		logEntry := APILogEntry{
			Timestamp:    startTime.Format("2006-01-02 15:04:05"),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Query:        queryParams,
			Headers:      headers,
			RequestBody:  requestBody,
			StatusCode:   c.Writer.Status(),
			ResponseBody: responseBody,
			Duration:     duration.String(),
			DurationMs:   duration.Milliseconds(),
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			UserID:       userID,
			UserEmail:    userEmail,
			Error:        errorMsg,
		}

		// === 11. WRITE TO LOG FILE (ASYNC) ===
		go writeLogToFile(logEntry, logsDir)

		// === 12. PRINT TO CONSOLE (DEVELOPMENT) ===
		if os.Getenv("APP_ENV") == "development" {
			printLogToConsole(logEntry)
		}
	}
}

// Write log ke file
func writeLogToFile(entry APILogEntry, logsDir string) {
	// Buat file log harian
	date := time.Now().Format("2006-01-02")
	filename := filepath.Join(logsDir, fmt.Sprintf("api_%s.log", date))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	// Write as pretty JSON
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v\n", err)
		return
	}

	// Write dengan separator
	file.WriteString("\n" + string(jsonData) + "\n")
	file.WriteString("========================================\n")
}

// Print ke console dengan warna
func printLogToConsole(entry APILogEntry) {
	statusColor := getStatusColor(entry.StatusCode)
	methodColor := getMethodColor(entry.Method)

	fmt.Printf("\n%s[%s]%s %s%s%s %s %s(%dms)%s",
		"\033[36m", entry.Timestamp, "\033[0m",
		methodColor, entry.Method, "\033[0m",
		entry.Path,
		statusColor, entry.DurationMs, "\033[0m",
	)

	if entry.StatusCode >= 400 {
		fmt.Printf(" %sStatus: %d%s", statusColor, entry.StatusCode, "\033[0m")
	}

	if entry.Error != "" {
		fmt.Printf(" %sError: %s%s", "\033[31m", entry.Error, "\033[0m")
	}

	if entry.UserID != nil {
		fmt.Printf(" %sUser: %v%s", "\033[35m", entry.UserID, "\033[0m")
	}

	fmt.Println()
}

// Mask auth token untuk keamanan
func maskAuthHeader(auth string) string {
	if auth == "" {
		return ""
	}
	if len(auth) > 20 {
		return auth[:7] + "..." + auth[len(auth)-4:]
	}
	return "Bearer ***"
}

// Warna untuk status code
func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[32m" // Green (Success)
	case status >= 300 && status < 400:
		return "\033[36m" // Cyan (Redirect)
	case status >= 400 && status < 500:
		return "\033[33m" // Yellow (Client Error)
	default:
		return "\033[31m" // Red (Server Error)
	}
}

// Warna untuk HTTP method
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Blue
	case "POST":
		return "\033[32m" // Green
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[35m" // Magenta
	default:
		return "\033[37m" // White
	}
}
