package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log" // Import the log package
	"net/http"
	"os"
	"time" // Import the time package
)

// OllamaGenerateRequest defines the structure for a request to the Ollama API.
type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaGenerateResponse defines the structure for a successful response from the Ollama API.
type OllamaGenerateResponse struct {
	Response string `json:"response"`
}

const commitMessagePrompt = "Based on the following 'git diff', generate a concise and semantic commit message that follows the Conventional Commits specification. The diff is:\n\n%s"

var ErrEmptyInput = errors.New("input from stdin is empty")

// run contains the core logic of the application.
func run(stdin io.Reader, stdout io.Writer, ollamaURL string, model string) error {
	log.Printf("Starting commit generation. Ollama URL: %s, Model: %s", ollamaURL, model)

	// --- 1. Read git diff from stdin ---
	diffBytes, err := io.ReadAll(stdin)
	if err != nil {
		log.Printf("ERROR: error reading from stdin: %v", err)
		return fmt.Errorf("error reading from stdin: %w", err)
	}
	if len(diffBytes) == 0 {
		log.Print("WARNING: Input from stdin is empty. No diff provided.")
		return ErrEmptyInput
	}
	log.Printf("Read %d bytes of Git diff from stdin.", len(diffBytes))

	// --- 2. Prepare the request for Ollama ---
	prompt := fmt.Sprintf(commitMessagePrompt, string(diffBytes))
	requestData := OllamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("ERROR: error marshalling JSON for Ollama: %v", err)
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	log.Print("JSON request for Ollama prepared.")

	// --- 3. Make the HTTP request with retries ---
	const maxRetries = 3
	var resp *http.Response
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		log.Printf("Making HTTP request to Ollama at: %s (Attempt %d/%d)", ollamaURL, i+1, maxRetries+1)

		// Create a new request body for each retry, as the reader gets consumed
		requestBody := bytes.NewBuffer(jsonData)
		
		resp, lastErr = http.Post(ollamaURL, "application/json", requestBody)
		if lastErr == nil && resp.StatusCode == http.StatusOK {
			log.Printf("HTTP request completed with status: %d", resp.StatusCode)
			break // Success!
		}

		if resp != nil { // Close the body if a response was received
			resp.Body.Close()
		}

		if i < maxRetries {
			// Decide if we should retry based on error or status code
			shouldRetry := true
			if lastErr == nil { // Got a response, but it was an error status
				// Don't retry on 4xx client errors
				if resp.StatusCode >= 400 && resp.StatusCode < 500 {
					shouldRetry = false
					log.Printf("Not retrying: Client error (status %d).", resp.StatusCode)
				}
			} else {
				// Retry on network errors
				log.Printf("Network error on attempt %d: %v", i+1, lastErr)
			}

			if shouldRetry {
				backoff := time.Duration(1<<uint(i)) * time.Second // 1s, 2s, 4s
				log.Printf("Retrying in %v...", backoff)
				time.Sleep(backoff)
				continue
			}
		}
		// If we are here, either max retries reached, or should not retry
		if lastErr != nil {
			log.Printf("ERROR: error making HTTP request to Ollama after %d attempts: %v", i+1, lastErr)
			return fmt.Errorf("error making request to Ollama after %d attempts: %w", i+1, lastErr)
		}
		// If no network error but status is not OK, this is the final error
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("ERROR: Ollama API returned error (status %d): %s", resp.StatusCode, string(body))
			return fmt.Errorf("error from Ollama API (status %d): %s", resp.StatusCode, string(body))
		}
	}
	// If the loop finishes without success, and resp is still nil, it means an unhandled error without a response.
	if resp == nil || resp.StatusCode != http.StatusOK {
		log.Print("ERROR: Could not get a successful response from Ollama after multiple attempts.")
		return fmt.Errorf("could not get a successful response from Ollama after multiple attempts")
	}
	defer resp.Body.Close() // Ensure body is closed after successful loop completion
	log.Print("Successful response from Ollama.")

	// --- 4. Decode the response and print the commit message ---
	var ollamaResp OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		log.Printf("ERROR: error decoding Ollama response: %v", err)
		return fmt.Errorf("error decoding Ollama response: %w", err)
	}
	log.Print("Ollama response decoded successfully.")


	_, err = fmt.Fprint(stdout, ollamaResp.Response)
	if err != nil {
		log.Printf("ERROR: error writing to stdout: %v", err)
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	log.Print("Commit message sent to stdout.")

	return nil
}

func main() {
	// --- Configuration ---
	ollamaURL := flag.String("ollama-url", "http://localhost:11434/api/generate", "Ollama API URL")
	model := flag.String("model", "qwen3-coder:30b", "Ollama model to use")
	flag.Parse()

	// Configure logger to write to stderr by default
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Add date, time, and file:line to log output

	if err := run(os.Stdin, os.Stdout, *ollamaURL, *model); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err) // This error is still printed to stderr for consistency with run's error logging
		os.Exit(1)
	}
	log.Print("Llamit CLI finished successfully.")
}

