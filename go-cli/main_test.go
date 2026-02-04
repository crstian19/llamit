package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Helper to capture log output
func captureLog(f func(), t *testing.T) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	// Restore default logger output at the end of the test
	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Restore default flags too
	}()
	log.SetFlags(0) // Disable flags for easier comparison in tests
	f()
	return buf.String()
}

func TestRun_IntegrationSuccess(t *testing.T) {
	// This is an integration test that makes a real network request.
	// It's skipped in -short mode.
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set your Ollama endpoint and model here to run integration tests
	// Example: ollamaURL := "http://localhost:11434/api/generate"
	ollamaURL := "http://localhost:11434/api/generate"
	model := "qwen3-coder:30b"

	// A sample git diff to send
	sampleDiff := `diff --git a/main.go b/main.go
index 1234567..abcdefg 100644
--- a/main.go
+++ b/main.go
@@ -1,5 +1,5 @@
 package main
 
 func main() {
-	println("Hello, World!")
+	println("Hello, Universe!")
 }
`
	stdin := strings.NewReader(sampleDiff)
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, ollamaURL, model, "conventional", "")
	}, t)

	if err != nil {
		t.Fatalf("run() returned an unexpected error: %v", err)
	}

	if stdout.Len() == 0 {
		t.Error("run() did not write anything to stdout, expected a commit message")
	}

	if !strings.Contains(logOutput, "Starting commit generation.") {
		t.Error("Expected log output to contain 'Starting commit generation.'")
	}
	if !strings.Contains(logOutput, "Successful response from Ollama.") {
		t.Error("Expected log output to contain 'Successful response from Ollama.'")
	}
	t.Logf("Integration test successful, received from Ollama: %s", stdout.String())
	t.Logf("Log output:\n%s", logOutput)
}

func TestRun_EmptyInput(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, "dummy-url", "dummy-model", "conventional", "")
	}, t)

	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("Expected error %v, got %v", ErrEmptyInput, err)
	}
	if !strings.Contains(logOutput, "WARNING: Input from stdin is empty. No diff provided.") {
		t.Error("Expected log output to contain 'WARNING: Input from stdin is empty. No diff provided.'")
	}
}

func TestRun_BadURL(t *testing.T) {
	sampleDiff := "diff --git a/file.txt b/file.txt"
	stdin := strings.NewReader(sampleDiff)
	stdout := new(bytes.Buffer)

	// Use a non-existent port on localhost for a guaranteed connection error
	badURL := "http://localhost:12345/api/generate"

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, badURL, "dummy-model", "conventional", "")
	}, t)

	if err == nil {
		t.Fatal("Expected an error for a bad URL, but got nil")
	}
	// Expect multiple retry attempts and a final error log
	if !strings.Contains(logOutput, "Making HTTP request to Ollama at:") {
		t.Error("Expected log output to contain 'Making HTTP request to Ollama at:'")
	}
	if !strings.Contains(logOutput, "ERROR: error making HTTP request to Ollama after") {
		t.Error("Expected log output to contain 'ERROR: error making HTTP request to Ollama after'")
	}
	t.Logf("Log output for BadURL:\n%s", logOutput)
}

func TestRun_RetrySuccess(t *testing.T) {
	// Simulate 2 failures then 1 success
	failCount := 2
	requestCounter := 0

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter++
		if requestCounter <= failCount {
			w.WriteHeader(http.StatusInternalServerError) // Simulate a transient error (e.g., 500)
			io.WriteString(w, `{"error": "Internal Server Error"}`)
			return
		}
		// Success on the third attempt
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: A new feature from retry"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, mockServer.URL, "dummy-model", "conventional", "")
	}, t)

	if err != nil {
		t.Fatalf("run() returned an unexpected error: %v", err)
	}

	if !strings.Contains(logOutput, "Retrying in ") {
		t.Error("Expected log output to contain retry messages")
	}
	if !strings.Contains(logOutput, "HTTP request completed with status: 200") {
		t.Error("Expected log output to show final success status 200")
	}
	if !strings.Contains(stdout.String(), "A new feature from retry") {
		t.Errorf("Expected stdout to contain success message, got: %s", stdout.String())
	}
	t.Logf("Log output for RetrySuccess:\n%s", logOutput)
}

func TestRun_RetryFailure(t *testing.T) {
	// Simulate always failing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Always fails"}`)
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, mockServer.URL, "dummy-model", "conventional", "")
	}, t)

	if err == nil {
		t.Fatal("Expected run() to return an error, but got nil")
	}
	t.Logf("Actual error from run: %v", err.Error())
	if !strings.Contains(err.Error(), `error from Ollama API (status 500): {"error": "Always fails"}`) {
		t.Errorf("Expected specific error about Ollama API status 500, got: %v", err)
	}

	if !strings.Contains(logOutput, "Retrying in ") {
		t.Error("Expected log output to contain retry messages")
	}
	if !strings.Contains(logOutput, `ERROR: Ollama API returned error (status 500): `) {
		t.Error("Expected log output to show final error after retries")
	}
	t.Logf("Log output for RetryFailure:\n%s", logOutput)
}

func TestRun_NoRetryOnClientError(t *testing.T) {
	requestCounter := 0
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter++
		w.WriteHeader(http.StatusBadRequest) // Simulate a client error (e.g., 400)
		io.WriteString(w, `{"error": "Bad Request"}`)
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, mockServer.URL, "dummy-model", "conventional", "")
	}, t)

	if err == nil {
		t.Fatal("Expected run() to return an error, but got nil")
	}
	t.Logf("Actual error from run: %v", err.Error())
	if !strings.Contains(err.Error(), `error from Ollama API (status 400): {"error": "Bad Request"}`) {
		t.Errorf("Expected specific error about Ollama API status 400, got: %v", err)
	}

	if requestCounter != 1 {
		t.Errorf("Expected only 1 request, but got %d (client error should not retry)", requestCounter)
	}
	if strings.Contains(logOutput, "Retrying in ") {
		t.Error("Did not expect log output to contain retry messages for client error")
	}
	t.Logf("Log output for NoRetryOnClientError:\n%s", logOutput)
}

func TestGetFormatTemplate(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		customTemplate string
		wantContains   string
	}{
		{"default", "", "", "Conventional Commits"},
		{"angular", "angular", "", "Angular"},
		{"gitmoji", "gitmoji", "", "Gitmoji"},
		{"karma", "karma", "", "Karma"},
		{"semantic", "semantic", "", "Semantic"},
		{"google", "google", "", "Google"},
		{"custom_provided", "custom", "My Custom Template", "My Custom Template"},
		{"custom_empty", "custom", "", "Conventional Commits"},
		{"invalid", "unknown", "", "Conventional Commits"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFormatTemplate(tt.format, tt.customTemplate)
			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("getFormatTemplate(%q, %q) = %q, want it to contain %q", tt.format, tt.customTemplate, got, tt.wantContains)
			}
		})
	}
}
