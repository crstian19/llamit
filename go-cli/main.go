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

// Format templates for different commit message styles
const (
	conventionalFormat = `Generate a concise commit message following the Conventional Commits format:
<type>(<scope>): <description>

[optional body]

Rules:
1. **Be extremely concise.** Use the fewest words possible.
2. First line must be 50 characters or less.
3. Use imperative mood (e.g., "add", "fix", "refactor").
4. Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore.
5. Scope is optional but recommended.
6. Summary must be lowercase and no period at the end.
7. Body is optional; if used, limit to 1-2 short bullet points about the "why".

The diff is:

%s`

	angularFormat = `Generate a concise commit message following the Angular format:
<type>(<scope>): <subject>

[optional body]

Rules:
1. **Brevity is key.** Keep it short and direct.
2. Subject line must be 50 characters or less.
3. Type must be: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert.
4. Scope is required.
5. Subject must be imperative, lowercase, no period.
6. Body is optional; wrap at 72 characters and keep it very brief.

The diff is:

%s`

	gitmojiFormat = `Generate a concise commit message following the Gitmoji format:
<emoji> <type>(<scope>): <description>

[optional body]

Rules:
1. **Be concise.** Focus on the main change.
2. Start with an appropriate gitmoji (✨, 🐛, 📝, 💄, ♻️, ⚡️, ✅, 🔧).
3. Keep first line under 50 characters (including emoji).
4. Use imperative mood.
5. Body is optional and should be very short.

The diff is:

%s`

	karmaFormat = `Generate a concise commit message following the Karma format:
<type>(<scope>): <subject>

[optional body]

Rules:
1. **Keep it short.** No unnecessary details.
2. Type: feat, fix, docs, style, refactor, perf, test, chore.
3. Subject must be imperative, present tense, no period.
4. First line must be 50 characters or less.
5. Body is optional and should be brief.

The diff is:

%s`

	semanticFormat = `Generate a concise commit message following the Semantic format:
<type>: <description>

[optional body]

Rules:
1. **Be brief and direct.**
2. Type: feat, fix, docs, style, refactor, perf, test, build, ops, chore.
3. First line should be 50 characters or less.
4. Use imperative mood.
5. Body is optional and should focus on "why".

The diff is:

%s`

	googleFormat = `Generate a concise commit message following the Google format:
<subject>

[optional body]

Rules:
1. **Maximum conciseness.**
2. Subject: concise summary in imperative mood, max 50 chars.
3. Body: explain the essential "why", keep it very short.
4. Wrap body at 72 characters.

The diff is:

%s`
)

var ErrEmptyInput = errors.New("input from stdin is empty")

// getFormatTemplate returns the appropriate format template based on the format name
func getFormatTemplate(format string, customTemplate string) string {
	switch format {
	case "angular":
		return angularFormat
	case "gitmoji":
		return gitmojiFormat
	case "karma":
		return karmaFormat
	case "semantic":
		return semanticFormat
	case "google":
		return googleFormat
	case "custom":
		if customTemplate != "" {
			return customTemplate + "\n\nThe diff is:\n\n%s"
		}
		// Fallback to conventional if custom is selected but no template provided
		return conventionalFormat
	case "conventional":
		fallthrough
	default:
		return conventionalFormat
	}
}

// run contains the core logic of the application.
func run(stdin io.Reader, stdout io.Writer, ollamaURL string, model string, format string, customTemplate string) error {
	log.Printf("Starting commit generation. Ollama URL: %s, Model: %s, Format: %s", ollamaURL, model, format)

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
	formatTemplate := getFormatTemplate(format, customTemplate)
	prompt := fmt.Sprintf(formatTemplate, string(diffBytes))
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
				if resp != nil { // Close it now that we decided to retry
					resp.Body.Close()
				}
				backoff := time.Duration(1<<uint(i)) * time.Second // 1s, 2s, 4s
				log.Printf("Retrying in %v...", backoff)
				time.Sleep(backoff)
				continue
			}
		}
		// If we are here, either max retries reached, or should not retry
		if lastErr != nil {
			if resp != nil {
				resp.Body.Close()
			}
			log.Printf("ERROR: error making HTTP request to Ollama after %d attempts: %v", i+1, lastErr)
			return fmt.Errorf("error making request to Ollama after %d attempts: %w", i+1, lastErr)
		}
		// If no network error but status is not OK, this is the final error
		if resp != nil {
			defer resp.Body.Close()
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
	format := flag.String("format", "conventional", "Commit message format (conventional, angular, gitmoji, karma, semantic, google, custom)")
	customTemplate := flag.String("custom-template", "", "Custom format template (only used when format is 'custom')")
	version := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *version {
		fmt.Println("Llamit CLI v0.2.1-format-fix")
		return
	}

	// Configure logger to write to stderr by default
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Add date, time, and file:line to log output

	if err := run(os.Stdin, os.Stdout, *ollamaURL, *model, *format, *customTemplate); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err) // This error is still printed to stderr for consistency with run's error logging
		os.Exit(1)
	}
	log.Print("Llamit CLI finished successfully.")
}
