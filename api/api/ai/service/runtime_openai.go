package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	appConfig "dodevops-api/common/config"
)

type aiRuntimeClient struct {
	enabled         bool
	provider        string
	baseURL         string
	apiKey          string
	model           string
	reasoningEffort string
	httpClient      *http.Client
}

type openAIResponsesRequest struct {
	Model           string                 `json:"model"`
	Instructions    string                 `json:"instructions,omitempty"`
	Input           []openAIInputMessage   `json:"input"`
	Store           bool                   `json:"store"`
	MaxOutputTokens int                    `json:"max_output_tokens,omitempty"`
	Reasoning       *openAIReasoningConfig `json:"reasoning,omitempty"`
	Stream          bool                   `json:"stream,omitempty"`
}

type openAIInputMessage struct {
	Role    string             `json:"role"`
	Content []openAIInputBlock `json:"content"`
}

type openAIInputBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type openAIReasoningConfig struct {
	Effort string `json:"effort,omitempty"`
}

type openAIResponsesResponse struct {
	Status string `json:"status,omitempty"`
	OutputText string `json:"output_text,omitempty"`
	Output []struct {
		Type    string `json:"type"`
		Text    string `json:"text,omitempty"`
		Content []struct {
			Type    string `json:"type"`
			Text    string `json:"text,omitempty"`
			Value   string `json:"value,omitempty"`
			Content []struct {
				Type  string `json:"type,omitempty"`
				Text  string `json:"text,omitempty"`
				Value string `json:"value,omitempty"`
			} `json:"content,omitempty"`
		} `json:"content"`
	} `json:"output"`
	Choices []struct {
		Text string `json:"text,omitempty"`
		Message struct {
			Role    string          `json:"role,omitempty"`
			Content json.RawMessage `json:"content,omitempty"`
		} `json:"message,omitempty"`
	} `json:"choices,omitempty"`
	Text json.RawMessage `json:"text,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

type aiRuntimeProbeResult struct {
	Reachable bool
	LastError string
	CheckedAt time.Time
}

const (
	aiRuntimeProbeTTL     = 30 * time.Second
	aiRuntimeProbeTimeout = 5 * time.Second
)

var (
	aiRuntimeProbeCacheMu sync.Mutex
	aiRuntimeProbeCache   = map[string]aiRuntimeProbeResult{}
	errAIRuntimeCompletedWithoutText = errors.New("ai runtime completed successfully but did not return any text payload")
)

func newAIRuntimeClient() *aiRuntimeClient {
	if appConfig.Config == nil {
		return nil
	}

	cfg := appConfig.Config.AI
	provider := strings.TrimSpace(strings.ToLower(cfg.Provider))
	if provider == "" {
		provider = "openai"
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = "gpt-5.4"
	}

	timeoutSeconds := cfg.TimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = 60
	}

	apiKey := strings.TrimSpace(cfg.APIKey)
	enabled := apiKey != ""
	if !enabled {
		return &aiRuntimeClient{
			enabled:  false,
			provider: provider,
			baseURL:  baseURL,
			model:    model,
		}
	}

	return &aiRuntimeClient{
		enabled:         true,
		provider:        provider,
		baseURL:         strings.TrimRight(baseURL, "/"),
		apiKey:          apiKey,
		model:           model,
		reasoningEffort: strings.TrimSpace(strings.ToLower(cfg.ReasoningEffort)),
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

func (c *aiRuntimeClient) IsEnabled() bool {
	return c != nil && c.enabled && c.provider == "openai"
}

func (c *aiRuntimeClient) Provider() string {
	if c == nil {
		return ""
	}
	return c.provider
}

func (c *aiRuntimeClient) Model() string {
	if c == nil {
		return ""
	}
	return c.model
}

func (c *aiRuntimeClient) Probe(ctx context.Context) aiRuntimeProbeResult {
	checkedAt := time.Now()
	if !c.IsEnabled() {
		return aiRuntimeProbeResult{CheckedAt: checkedAt}
	}

	cacheKey := c.probeCacheKey()
	if cached, ok := loadAIRuntimeProbeCache(cacheKey, checkedAt); ok {
		return cached
	}

	if ctx == nil {
		ctx = context.Background()
	}
	probeTimeout := aiRuntimeProbeTimeout
	if c.httpClient != nil && c.httpClient.Timeout > 0 && c.httpClient.Timeout < probeTimeout {
		probeTimeout = c.httpClient.Timeout
	}
	probeCtx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	_, err := c.createTextResponseForProbe(probeCtx, "Reply with OK only.", "health check", 8)
	result := aiRuntimeProbeResult{
		Reachable: err == nil,
		CheckedAt: checkedAt,
	}
	if err != nil {
		result.LastError = normalizeAIRuntimeProbeError(err)
	}
	storeAIRuntimeProbeCache(cacheKey, result)
	return result
}

func (c *aiRuntimeClient) CreateTextResponse(ctx context.Context, instructions, input string, maxOutputTokens int) (string, error) {
	if !c.IsEnabled() {
		return "", fmt.Errorf("ai runtime is not enabled")
	}

	reqBody := c.buildTextResponseRequest(instructions, input, maxOutputTokens)

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp openAIResponsesResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != nil && strings.TrimSpace(errResp.Error.Message) != "" {
			return "", fmt.Errorf("openai responses error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("openai responses request failed: %s", strings.TrimSpace(string(body)))
	}

	var response openAIResponsesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	text := extractTextFromAIRuntimeResponse(response)
	if text == "" {
		if response.Error == nil && strings.EqualFold(strings.TrimSpace(response.Status), "completed") {
			streamText, streamErr := c.createTextResponseFromResponsesStream(ctx, reqBody)
			if streamErr == nil && strings.TrimSpace(streamText) != "" {
				return streamText, nil
			}
			if streamErr != nil {
				return "", streamErr
			}
			return "", errAIRuntimeCompletedWithoutText
		}
		return "", fmt.Errorf("openai responses returned empty text")
	}
	return text, nil
}

func (c *aiRuntimeClient) createTextResponseForProbe(ctx context.Context, instructions, input string, maxOutputTokens int) (string, error) {
	reqBody := c.buildTextResponseRequestForProbe(instructions, input, maxOutputTokens)

	streamText, streamErr := c.createTextResponseFromResponsesStream(ctx, reqBody)
	if streamErr == nil && strings.TrimSpace(streamText) != "" {
		return streamText, nil
	}

	return c.CreateTextResponse(ctx, instructions, input, maxOutputTokens)
}

func (c *aiRuntimeClient) probeCacheKey() string {
	if c == nil {
		return ""
	}
	return strings.Join([]string{
		c.provider,
		c.baseURL,
		c.model,
		c.reasoningEffort,
	}, "|")
}

func loadAIRuntimeProbeCache(cacheKey string, now time.Time) (aiRuntimeProbeResult, bool) {
	if strings.TrimSpace(cacheKey) == "" {
		return aiRuntimeProbeResult{}, false
	}

	aiRuntimeProbeCacheMu.Lock()
	defer aiRuntimeProbeCacheMu.Unlock()

	result, ok := aiRuntimeProbeCache[cacheKey]
	if !ok {
		return aiRuntimeProbeResult{}, false
	}
	if now.Sub(result.CheckedAt) > aiRuntimeProbeTTL {
		delete(aiRuntimeProbeCache, cacheKey)
		return aiRuntimeProbeResult{}, false
	}
	return result, true
}

func storeAIRuntimeProbeCache(cacheKey string, result aiRuntimeProbeResult) {
	if strings.TrimSpace(cacheKey) == "" {
		return
	}

	aiRuntimeProbeCacheMu.Lock()
	defer aiRuntimeProbeCacheMu.Unlock()
	aiRuntimeProbeCache[cacheKey] = result
}

func normalizeAIRuntimeProbeError(err error) string {
	if err == nil {
		return ""
	}
	message := strings.Join(strings.Fields(strings.TrimSpace(err.Error())), " ")
	lowerMessage := strings.ToLower(message)
	switch {
	case errors.Is(err, errAIRuntimeCompletedWithoutText),
		strings.Contains(lowerMessage, "completed successfully but did not return any text payload"),
		strings.Contains(lowerMessage, "returned empty text"):
		return "模型网关已响应，但当前返回里没有标准文本内容，系统将先回退到内置逻辑。"
	case strings.Contains(lowerMessage, "context deadline exceeded"),
		strings.Contains(lowerMessage, "client.timeout exceeded"),
		strings.Contains(lowerMessage, "timeout awaiting response headers"):
		return "模型网关探测超时，当前先回退到内置逻辑，请稍后重试。"
	case strings.Contains(lowerMessage, "invalid character"),
		strings.Contains(lowerMessage, "cannot unmarshal"),
		strings.Contains(lowerMessage, "unexpected end of json input"):
		return "模型网关已返回结果，但返回格式和当前解析规则不兼容，系统将先回退到内置逻辑。"
	}
	const maxLen = 240
	if len(message) > maxLen {
		return message[:maxLen] + "..."
	}
	return message
}

func extractTextFromAIRuntimeResponse(response openAIResponsesResponse) string {
	var texts []string

	appendAIRuntimeText(&texts, response.OutputText)
	for _, item := range response.Output {
		appendAIRuntimeText(&texts, item.Text)
		for _, content := range item.Content {
			appendAIRuntimeText(&texts, content.Text)
			appendAIRuntimeText(&texts, content.Value)
			for _, nested := range content.Content {
				appendAIRuntimeText(&texts, nested.Text)
				appendAIRuntimeText(&texts, nested.Value)
			}
		}
	}
	for _, choice := range response.Choices {
		appendAIRuntimeText(&texts, choice.Text)
		appendAIRuntimeText(&texts, extractTextFromChoiceMessageContent(choice.Message.Content))
	}
	appendAIRuntimeText(&texts, extractTextFromRawField(response.Text))

	return strings.TrimSpace(strings.Join(texts, "\n"))
}

func appendAIRuntimeText(target *[]string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	*target = append(*target, value)
}

func extractTextFromChoiceMessageContent(raw json.RawMessage) string {
	if len(bytes.TrimSpace(raw)) == 0 {
		return ""
	}

	var plain string
	if err := json.Unmarshal(raw, &plain); err == nil {
		return strings.TrimSpace(plain)
	}

	var object struct {
		Text    string `json:"text,omitempty"`
		Value   string `json:"value,omitempty"`
		Content string `json:"content,omitempty"`
	}
	if err := json.Unmarshal(raw, &object); err == nil {
		return firstNonEmpty(strings.TrimSpace(object.Text), strings.TrimSpace(object.Value), strings.TrimSpace(object.Content))
	}

	var blocks []struct {
		Text    string `json:"text,omitempty"`
		Value   string `json:"value,omitempty"`
		Content string `json:"content,omitempty"`
	}
	if err := json.Unmarshal(raw, &blocks); err == nil {
		var texts []string
		for _, block := range blocks {
			appendAIRuntimeText(&texts, block.Text)
			appendAIRuntimeText(&texts, block.Value)
			appendAIRuntimeText(&texts, block.Content)
		}
		return strings.TrimSpace(strings.Join(texts, "\n"))
	}

	return ""
}

func extractTextFromRawField(raw json.RawMessage) string {
	if len(bytes.TrimSpace(raw)) == 0 {
		return ""
	}

	var plain string
	if err := json.Unmarshal(raw, &plain); err == nil {
		return strings.TrimSpace(plain)
	}

	var object struct {
		Value   string `json:"value,omitempty"`
		Content string `json:"content,omitempty"`
		Text    string `json:"text,omitempty"`
	}
	if err := json.Unmarshal(raw, &object); err == nil {
		return firstNonEmpty(strings.TrimSpace(object.Value), strings.TrimSpace(object.Content), strings.TrimSpace(object.Text))
	}

	return ""
}

func (c *aiRuntimeClient) buildTextResponseRequest(instructions, input string, maxOutputTokens int) openAIResponsesRequest {
	return c.buildTextResponseRequestWithReasoning(instructions, input, maxOutputTokens, true)
}

func (c *aiRuntimeClient) buildTextResponseRequestForProbe(instructions, input string, maxOutputTokens int) openAIResponsesRequest {
	return c.buildTextResponseRequestWithReasoning(instructions, input, maxOutputTokens, false)
}

func (c *aiRuntimeClient) buildTextResponseRequestWithReasoning(instructions, input string, maxOutputTokens int, includeReasoning bool) openAIResponsesRequest {
	reqBody := openAIResponsesRequest{
		Model:        c.model,
		Instructions: strings.TrimSpace(instructions),
		Input: []openAIInputMessage{
			{
				Role: "user",
				Content: []openAIInputBlock{
					{
						Type: "input_text",
						Text: input,
					},
				},
			},
		},
		Store: false,
	}
	if maxOutputTokens > 0 {
		reqBody.MaxOutputTokens = maxOutputTokens
	}
	if includeReasoning && c.reasoningEffort != "" {
		reqBody.Reasoning = &openAIReasoningConfig{Effort: c.reasoningEffort}
	}
	return reqBody
}

func (c *aiRuntimeClient) createTextResponseFromResponsesStream(ctx context.Context, reqBody openAIResponsesRequest) (string, error) {
	if c == nil || c.httpClient == nil {
		return "", errAIRuntimeCompletedWithoutText
	}

	reqBody.Stream = true
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return "", readErr
		}
		return "", fmt.Errorf("openai responses stream request failed: %s", strings.TrimSpace(string(body)))
	}

	text, err := extractTextFromResponsesEventStream(resp.Body)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(text) == "" {
		return "", errAIRuntimeCompletedWithoutText
	}
	return text, nil
}

func extractTextFromResponsesEventStream(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var deltaBuilder strings.Builder
	finalText := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" || data == "[DONE]" {
			continue
		}

		eventType, chunkText, err := extractTextFromResponsesStreamChunk([]byte(data))
		if err != nil {
			continue
		}
		switch eventType {
		case "response.output_text.delta":
			deltaBuilder.WriteString(chunkText)
		case "response.output_text.done", "response.content_part.done", "response.completed":
			if deltaBuilder.Len() > 0 {
				return strings.TrimSpace(deltaBuilder.String()), nil
			}
			if strings.TrimSpace(chunkText) != "" {
				return strings.TrimSpace(chunkText), nil
			}
			if strings.TrimSpace(finalText) == "" {
				finalText = chunkText
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if deltaBuilder.Len() > 0 {
		return strings.TrimSpace(deltaBuilder.String()), nil
	}
	return strings.TrimSpace(finalText), nil
}

func extractTextFromResponsesStreamChunk(data []byte) (string, string, error) {
	var envelope struct {
		Type     string                 `json:"type,omitempty"`
		Delta    string                 `json:"delta,omitempty"`
		Text     string                 `json:"text,omitempty"`
		Response openAIResponsesResponse `json:"response,omitempty"`
		Part     struct {
			Type string `json:"type,omitempty"`
			Text string `json:"text,omitempty"`
		} `json:"part,omitempty"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return "", "", err
	}

	switch envelope.Type {
	case "response.output_text.delta":
		return envelope.Type, strings.TrimSpace(envelope.Delta), nil
	case "response.output_text.done":
		return envelope.Type, strings.TrimSpace(envelope.Text), nil
	case "response.content_part.done":
		if strings.EqualFold(strings.TrimSpace(envelope.Part.Type), "output_text") {
			return envelope.Type, strings.TrimSpace(envelope.Part.Text), nil
		}
	case "response.completed":
		return envelope.Type, extractTextFromAIRuntimeResponse(envelope.Response), nil
	}

	return envelope.Type, "", nil
}
