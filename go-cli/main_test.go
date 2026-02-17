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

func captureLog(f func(), t *testing.T) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}()
	log.SetFlags(0)
	f()
	return buf.String()
}

func TestRun_IntegrationSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ollamaURL := "http://localhost:11434/api/generate"
	model := "qwen2.5-coder:7b"

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
		err = run(stdin, stdout, RunConfig{
			ollamaURL: ollamaURL,
			model:     model,
			format:    "conventional",
		})
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
		err = run(stdin, stdout, RunConfig{
			ollamaURL: "dummy-url",
			model:     "dummy-model",
			format:    "conventional",
		})
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

	badURL := "http://localhost:12345/api/generate"

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, RunConfig{
			ollamaURL: badURL,
			model:     "dummy-model",
			format:    "conventional",
		})
	}, t)

	if err == nil {
		t.Fatal("Expected an error for a bad URL, but got nil")
	}
	if !strings.Contains(logOutput, "Making HTTP request to Ollama at:") {
		t.Error("Expected log output to contain 'Making HTTP request to Ollama at:'")
	}
	if !strings.Contains(logOutput, "ERROR: error making HTTP request to Ollama after") {
		t.Error("Expected log output to contain 'ERROR: error making HTTP request to Ollama after'")
	}
	t.Logf("Log output for BadURL:\n%s", logOutput)
}

func TestRun_RetrySuccess(t *testing.T) {
	failCount := 2
	requestCounter := 0

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter++
		if requestCounter <= failCount {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"error": "Internal Server Error"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: A new feature from retry"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, RunConfig{
			ollamaURL: mockServer.URL,
			model:     "dummy-model",
			format:    "conventional",
		})
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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Always fails"}`)
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, RunConfig{
			ollamaURL: mockServer.URL,
			model:     "dummy-model",
			format:    "conventional",
		})
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
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Bad Request"}`)
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	var err error
	logOutput := captureLog(func() {
		err = run(stdin, stdout, RunConfig{
			ollamaURL: mockServer.URL,
			model:     "dummy-model",
			format:    "conventional",
		})
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

func TestRun_WithKeepAlive(t *testing.T) {
	var receivedReq OllamaGenerateRequest

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: test"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	err := run(stdin, stdout, RunConfig{
		ollamaURL: mockServer.URL,
		model:     "test-model",
		format:    "conventional",
		keepAlive: "5m",
	})

	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if receivedReq.KeepAlive != "5m" {
		t.Errorf("Expected keep_alive='5m', got '%s'", receivedReq.KeepAlive)
	}
}

func TestRun_WithTemperature(t *testing.T) {
	var receivedReq OllamaGenerateRequest

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: test"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	temp := 0.7
	err := run(stdin, stdout, RunConfig{
		ollamaURL: mockServer.URL,
		model:     "test-model",
		format:    "conventional",
		options: &OllamaOptions{
			Temperature: &temp,
		},
	})

	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if receivedReq.Options == nil {
		t.Fatal("Expected options to be present in request")
	}
	if receivedReq.Options["temperature"] != 0.7 {
		t.Errorf("Expected temperature=0.7, got %v", receivedReq.Options["temperature"])
	}
}

func TestRun_WithMultipleOptions(t *testing.T) {
	var receivedReq OllamaGenerateRequest

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: test"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	temp := 0.9
	topK := 40
	topP := 0.95
	err := run(stdin, stdout, RunConfig{
		ollamaURL: mockServer.URL,
		model:     "test-model",
		format:    "conventional",
		keepAlive: "10m",
		options: &OllamaOptions{
			Temperature: &temp,
			TopK:        &topK,
			TopP:        &topP,
		},
	})

	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if receivedReq.KeepAlive != "10m" {
		t.Errorf("Expected keep_alive='10m', got '%s'", receivedReq.KeepAlive)
	}
	if receivedReq.Options["temperature"] != 0.9 {
		t.Errorf("Expected temperature=0.9, got %v", receivedReq.Options["temperature"])
	}
	if receivedReq.Options["top_k"].(float64) != 40 {
		t.Errorf("Expected top_k=40, got %v", receivedReq.Options["top_k"])
	}
	if receivedReq.Options["top_p"] != 0.95 {
		t.Errorf("Expected top_p=0.95, got %v", receivedReq.Options["top_p"])
	}
}

func TestRun_NoOptionsWhenNotSet(t *testing.T) {
	var receivedReq OllamaGenerateRequest

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: test"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	err := run(stdin, stdout, RunConfig{
		ollamaURL: mockServer.URL,
		model:     "test-model",
		format:    "conventional",
		options:   &OllamaOptions{},
	})

	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	if receivedReq.KeepAlive != "" {
		t.Errorf("Expected keep_alive to be empty, got '%s'", receivedReq.KeepAlive)
	}
	if receivedReq.Options != nil && len(receivedReq.Options) > 0 {
		t.Errorf("Expected options to be nil or empty, got %v", receivedReq.Options)
	}
}

func TestRun_WithStopSequences(t *testing.T) {
	var receivedReq map[string]interface{}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OllamaGenerateResponse{Response: "feat: test"})
	}))
	defer mockServer.Close()

	stdin := strings.NewReader("diff --git a/file.txt b/file.txt")
	stdout := new(bytes.Buffer)

	err := run(stdin, stdout, RunConfig{
		ollamaURL: mockServer.URL,
		model:     "test-model",
		format:    "conventional",
		options: &OllamaOptions{
			Stop: []string{"STOP1", "STOP2"},
		},
	})

	if err != nil {
		t.Fatalf("run() returned error: %v", err)
	}

	options, ok := receivedReq["options"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected options to be present in request")
	}
	stop, ok := options["stop"].([]interface{})
	if !ok {
		t.Fatal("Expected stop to be an array")
	}
	if len(stop) != 2 {
		t.Errorf("Expected 2 stop sequences, got %d", len(stop))
	}
	if stop[0] != "STOP1" || stop[1] != "STOP2" {
		t.Errorf("Expected stop=['STOP1', 'STOP2'], got %v", stop)
	}
}

func TestOllamaOptions_ToMap(t *testing.T) {
	temp := 0.5
	topK := 20
	stop := []string{"end"}

	options := &OllamaOptions{
		Temperature: &temp,
		TopK:        &topK,
		Stop:        stop,
	}

	result := options.ToMap()

	if len(result) != 3 {
		t.Errorf("Expected 3 options, got %d", len(result))
	}
	if result["temperature"] != 0.5 {
		t.Errorf("Expected temperature=0.5, got %v", result["temperature"])
	}
	if result["top_k"] != 20 {
		t.Errorf("Expected top_k=20, got %v", result["top_k"])
	}
}

func TestOllamaOptions_HasOptions(t *testing.T) {
	empty := &OllamaOptions{}
	if empty.HasOptions() {
		t.Error("Expected empty options to return false")
	}

	temp := 0.5
	withTemp := &OllamaOptions{Temperature: &temp}
	if !withTemp.HasOptions() {
		t.Error("Expected options with temperature to return true")
	}
}

func TestParseStopString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"STOP1", []string{"STOP1"}},
		{"STOP1,STOP2", []string{"STOP1", "STOP2"}},
		{" STOP1 , STOP2 ", []string{"STOP1", "STOP2"}},
		{",STOP1,,STOP2,", []string{"STOP1", "STOP2"}},
	}

	for _, tt := range tests {
		result := parseStopString(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseStopString(%q) = %v, expected %v", tt.input, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("parseStopString(%q)[%d] = %q, expected %q", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestGetFloatPtrIfSet(t *testing.T) {
	if getFloatPtrIfSet(0) != nil {
		t.Error("Expected nil for 0 value")
	}
	if getFloatPtrIfSet(1) == nil || *getFloatPtrIfSet(1) != 1 {
		t.Error("Expected pointer to 1 for value 1")
	}
}

func TestGetIntPtrIfSet(t *testing.T) {
	if getIntPtrIfSet(0) != nil {
		t.Error("Expected nil for 0 value")
	}
	if getIntPtrIfSet(1) == nil || *getIntPtrIfSet(1) != 1 {
		t.Error("Expected pointer to 1 for value 1")
	}
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
