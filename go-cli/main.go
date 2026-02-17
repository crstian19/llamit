package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type OllamaGenerateRequest struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	Stream    bool                   `json:"stream"`
	KeepAlive string                 `json:"keep_alive,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

type OllamaGenerateResponse struct {
	Response string `json:"response"`
}

type OllamaOptions struct {
	Temperature   *float64 `json:"temperature,omitempty"`
	TopK          *int     `json:"top_k,omitempty"`
	TopP          *float64 `json:"top_p,omitempty"`
	NumCtx        *int     `json:"num_ctx,omitempty"`
	NumPredict    *int     `json:"num_predict,omitempty"`
	RepeatPenalty *float64 `json:"repeat_penalty,omitempty"`
	RepeatLastN   *int     `json:"repeat_last_n,omitempty"`
	Seed          *int     `json:"seed,omitempty"`
	NumGPU        *int     `json:"num_gpu,omitempty"`
	NumThread     *int     `json:"num_thread,omitempty"`
	MinP          *float64 `json:"min_p,omitempty"`
	TfsZ          *float64 `json:"tfs_z,omitempty"`
	Mirostat      *int     `json:"mirostat,omitempty"`
	MirostatEta   *float64 `json:"mirostat_eta,omitempty"`
	MirostatTau   *float64 `json:"mirostat_tau,omitempty"`
	Stop          []string `json:"stop,omitempty"`
}

func (o *OllamaOptions) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	if o.Temperature != nil {
		result["temperature"] = *o.Temperature
	}
	if o.TopK != nil {
		result["top_k"] = *o.TopK
	}
	if o.TopP != nil {
		result["top_p"] = *o.TopP
	}
	if o.NumCtx != nil {
		result["num_ctx"] = *o.NumCtx
	}
	if o.NumPredict != nil {
		result["num_predict"] = *o.NumPredict
	}
	if o.RepeatPenalty != nil {
		result["repeat_penalty"] = *o.RepeatPenalty
	}
	if o.RepeatLastN != nil {
		result["repeat_last_n"] = *o.RepeatLastN
	}
	if o.Seed != nil {
		result["seed"] = *o.Seed
	}
	if o.NumGPU != nil {
		result["num_gpu"] = *o.NumGPU
	}
	if o.NumThread != nil {
		result["num_thread"] = *o.NumThread
	}
	if o.MinP != nil {
		result["min_p"] = *o.MinP
	}
	if o.TfsZ != nil {
		result["tfs_z"] = *o.TfsZ
	}
	if o.Mirostat != nil {
		result["mirostat"] = *o.Mirostat
	}
	if o.MirostatEta != nil {
		result["mirostat_eta"] = *o.MirostatEta
	}
	if o.MirostatTau != nil {
		result["mirostat_tau"] = *o.MirostatTau
	}
	if len(o.Stop) > 0 {
		result["stop"] = o.Stop
	}

	return result
}

func (o *OllamaOptions) HasOptions() bool {
	return o.Temperature != nil || o.TopK != nil || o.TopP != nil ||
		o.NumCtx != nil || o.NumPredict != nil || o.RepeatPenalty != nil ||
		o.RepeatLastN != nil || o.Seed != nil || o.NumGPU != nil ||
		o.NumThread != nil || o.MinP != nil || o.TfsZ != nil ||
		o.Mirostat != nil || o.MirostatEta != nil || o.MirostatTau != nil ||
		len(o.Stop) > 0
}

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
8. **IMPORTANT: DO NOT use markdown.** No backticks, no code blocks. Output plain text only.

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
7. **IMPORTANT: DO NOT use markdown.** No backticks, no code blocks. Output plain text only.

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
6. **IMPORTANT: DO NOT use markdown.** No backticks, no code blocks. Output plain text only.

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
6. **IMPORTANT: DO NOT use markdown.** No backticks, no code blocks. Output plain text only.

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
6. **IMPORTANT: DO NOT use markdown.** No backticks, no code blocks. Output plain text only.

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
5. **IMPORTANT: DO NOT use markdown.** No backticks, no code blocks. Output plain text only.

The diff is:

%s`
)

var ErrEmptyInput = errors.New("input from stdin is empty")

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
		return conventionalFormat
	case "conventional":
		fallthrough
	default:
		return conventionalFormat
	}
}

type RunConfig struct {
	ollamaURL      string
	model          string
	format         string
	customTemplate string
	keepAlive      string
	options        *OllamaOptions
}

func run(stdin io.Reader, stdout io.Writer, config RunConfig) error {
	log.Printf("Starting commit generation. Ollama URL: %s, Model: %s, Format: %s", config.ollamaURL, config.model, config.format)

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

	formatTemplate := getFormatTemplate(config.format, config.customTemplate)
	prompt := fmt.Sprintf(formatTemplate, string(diffBytes))

	requestData := OllamaGenerateRequest{
		Model:     config.model,
		Prompt:    prompt,
		Stream:    false,
		KeepAlive: config.keepAlive,
	}

	if config.options != nil && config.options.HasOptions() {
		requestData.Options = config.options.ToMap()
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("ERROR: error marshalling JSON for Ollama: %v", err)
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	log.Print("JSON request for Ollama prepared.")

	const maxRetries = 3
	var resp *http.Response
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		log.Printf("Making HTTP request to Ollama at: %s (Attempt %d/%d)", config.ollamaURL, i+1, maxRetries+1)

		requestBody := bytes.NewBuffer(jsonData)

		resp, lastErr = http.Post(config.ollamaURL, "application/json", requestBody)
		if lastErr == nil && resp.StatusCode == http.StatusOK {
			log.Printf("HTTP request completed with status: %d", resp.StatusCode)
			break
		}

		if i < maxRetries {
			shouldRetry := true
			if lastErr == nil {
				if resp.StatusCode >= 400 && resp.StatusCode < 500 {
					shouldRetry = false
					log.Printf("Not retrying: Client error (status %d).", resp.StatusCode)
				}
			} else {
				log.Printf("Network error on attempt %d: %v", i+1, lastErr)
			}

			if shouldRetry {
				if resp != nil {
					resp.Body.Close()
				}
				backoff := time.Duration(1<<uint(i)) * time.Second
				log.Printf("Retrying in %v...", backoff)
				time.Sleep(backoff)
				continue
			}
		}
		if lastErr != nil {
			if resp != nil {
				resp.Body.Close()
			}
			log.Printf("ERROR: error making HTTP request to Ollama after %d attempts: %v", i+1, lastErr)
			return fmt.Errorf("error making request to Ollama after %d attempts: %w", i+1, lastErr)
		}
		if resp != nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			log.Printf("ERROR: Ollama API returned error (status %d): %s", resp.StatusCode, string(body))
			return fmt.Errorf("error from Ollama API (status %d): %s", resp.StatusCode, string(body))
		}
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		log.Print("ERROR: Could not get a successful response from Ollama after multiple attempts.")
		return fmt.Errorf("could not get a successful response from Ollama after multiple attempts")
	}
	defer resp.Body.Close()
	log.Print("Successful response from Ollama.")

	var ollamaResp OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		log.Printf("ERROR: error decoding Ollama response: %v", err)
		return fmt.Errorf("error decoding Ollama response: %w", err)
	}
	log.Print("Ollama response decoded successfully.")

	commitMsg := cleanResponse(ollamaResp.Response)

	_, err = fmt.Fprint(stdout, commitMsg)
	if err != nil {
		log.Printf("ERROR: error writing to stdout: %v", err)
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	log.Print("Commit message sent to stdout.")

	return nil
}

func cleanResponse(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		lines := strings.Split(s, "\n")
		if len(lines) > 1 && strings.HasPrefix(lines[0], "```") {
			lastIdx := -1
			for i := len(lines) - 1; i > 0; i-- {
				if strings.TrimSpace(lines[i]) == "```" {
					lastIdx = i
					break
				}
			}
			if lastIdx != -1 {
				s = strings.Join(lines[1:lastIdx], "\n")
			} else {
				s = strings.Join(lines[1:], "\n")
			}
		}
	}
	s = strings.ReplaceAll(s, "```", "")
	s = strings.ReplaceAll(s, "`", "")
	return strings.TrimSpace(s)
}

func main() {
	ollamaURL := flag.String("ollama-url", "http://localhost:11434/api/generate", "Ollama API URL")
	model := flag.String("model", "qwen2.5-coder:7b", "Ollama model to use")
	format := flag.String("format", "conventional", "Commit message format (conventional, angular, gitmoji, karma, semantic, google, custom)")
	customTemplate := flag.String("custom-template", "", "Custom format template (only used when format is 'custom')")
	keepAlive := flag.String("keep-alive", "", "How long to keep the model loaded in memory (e.g., '5m', '60', '-1' for always)")
	version := flag.Bool("version", false, "Print version and exit")

	temperature := flag.Float64("temperature", 0, "Temperature for the model (0.0-2.0)")
	topK := flag.Int("top-k", 0, "Top-k sampling parameter")
	topP := flag.Float64("top-p", 0, "Top-p sampling parameter")
	numCtx := flag.Int("num-ctx", 0, "Context window size")
	numPredict := flag.Int("num-predict", 0, "Maximum tokens to predict")
	repeatPenalty := flag.Float64("repeat-penalty", 0, "Repetition penalty")
	repeatLastN := flag.Int("repeat-last-n", 0, "Look back distance for repetition penalty")
	seed := flag.Int("seed", 0, "Random seed (-1 for random)")
	numGPU := flag.Int("num-gpu", 0, "Number of GPUs to use")
	numThread := flag.Int("num-thread", 0, "Number of threads to use")
	minP := flag.Float64("min-p", 0, "Minimum probability sampling")
	tfsZ := flag.Float64("tfs-z", 0, "Tail free sampling parameter")
	mirostat := flag.Int("mirostat", 0, "Mirostat sampling (0, 1, 2)")
	mirostatEta := flag.Float64("mirostat-eta", 0, "Mirostat eta parameter")
	mirostatTau := flag.Float64("mirostat-tau", 0, "Mirostat tau parameter")
	stop := flag.String("stop", "", "Stop sequences (comma-separated)")

	flag.Parse()

	if *version {
		fmt.Println("Llamit CLI v0.3.0-ollama-options")
		return
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	options := &OllamaOptions{
		Temperature:   getFloatPtrIfSet(*temperature),
		TopK:          getIntPtrIfSet(*topK),
		TopP:          getFloatPtrIfSet(*topP),
		NumCtx:        getIntPtrIfSet(*numCtx),
		NumPredict:    getIntPtrIfSet(*numPredict),
		RepeatPenalty: getFloatPtrIfSet(*repeatPenalty),
		RepeatLastN:   getIntPtrIfSet(*repeatLastN),
		Seed:          getIntPtrIfSet(*seed),
		NumGPU:        getIntPtrIfSet(*numGPU),
		NumThread:     getIntPtrIfSet(*numThread),
		MinP:          getFloatPtrIfSet(*minP),
		TfsZ:          getFloatPtrIfSet(*tfsZ),
		Mirostat:      getIntPtrIfSet(*mirostat),
		MirostatEta:   getFloatPtrIfSet(*mirostatEta),
		MirostatTau:   getFloatPtrIfSet(*mirostatTau),
		Stop:          parseStopString(*stop),
	}

	config := RunConfig{
		ollamaURL:      *ollamaURL,
		model:          *model,
		format:         *format,
		customTemplate: *customTemplate,
		keepAlive:      *keepAlive,
		options:        options,
	}

	if err := run(os.Stdin, os.Stdout, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	log.Print("Llamit CLI finished successfully.")
}

func getFloatPtrIfSet(value float64) *float64 {
	if value != 0 {
		return &value
	}
	return nil
}

func getIntPtrIfSet(value int) *int {
	if value != 0 {
		return &value
	}
	return nil
}

func parseStopString(stopStr string) []string {
	if stopStr == "" {
		return nil
	}
	parts := strings.Split(stopStr, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
