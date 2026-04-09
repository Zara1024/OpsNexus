package service

import (
	"encoding/json"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestAIRuntimeClient(t *testing.T, handler http.HandlerFunc) *aiRuntimeClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &aiRuntimeClient{
		enabled:    true,
		provider:   "openai",
		baseURL:    server.URL,
		apiKey:     "test-key",
		model:      "gpt-5.4",
		httpClient: server.Client(),
	}
}

func TestCreateTextResponseSupportsOutputTextCompatibilityField(t *testing.T) {
	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_compat",
			"status": "completed",
			"output": [],
			"output_text": "OK"
		}`))
	})

	got, err := client.CreateTextResponse(context.Background(), "Reply with OK only.", "health check", 8)
	if err != nil {
		t.Fatalf("expected compatibility output_text to be accepted, got error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("expected output_text content, got %q", got)
	}
}

func TestCreateTextResponseSupportsChatCompletionChoicesContent(t *testing.T) {
	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "chatcmpl_compat",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "OK"
					}
				}
			]
		}`))
	})

	got, err := client.CreateTextResponse(context.Background(), "Reply with OK only.", "health check", 8)
	if err != nil {
		t.Fatalf("expected chat completions compatibility content to be accepted, got error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("expected choices message content, got %q", got)
	}
}

func TestCreateTextResponseReturnsCompatibilityErrorWhenGatewayCompletesWithoutText(t *testing.T) {
	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_beecode",
			"status": "completed",
			"error": null,
			"output": [],
			"text": {
				"format": {
					"type": "text"
				},
				"verbosity": "medium"
			}
		}`))
	})

	_, err := client.CreateTextResponse(context.Background(), "Reply with OK only.", "health check", 8)
	if err == nil {
		t.Fatalf("expected compatibility error for empty completed response")
	}
	if got, want := err.Error(), "ai runtime completed successfully but did not return any text payload"; got != want {
		t.Fatalf("expected compatibility error %q, got %q", want, got)
	}
}

func TestCreateTextResponseFallsBackToResponsesStreamWhenCompletedPayloadHasNoText(t *testing.T) {
	requestCount := 0
	streamSeen := false

	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("expected valid json payload, got error: %v", err)
		}

		if payload["stream"] == true {
			streamSeen = true
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("event: response.output_text.delta\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"OK\"}\n\n"))
			_, _ = w.Write([]byte("event: response.output_text.done\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.done\",\"text\":\"OK\"}\n\n"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_beecode",
			"status": "completed",
			"error": null,
			"output": []
		}`))
	})

	got, err := client.CreateTextResponse(context.Background(), "Reply with OK only.", "health check", 8)
	if err != nil {
		t.Fatalf("expected stream fallback to recover text, got error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("expected stream fallback text %q, got %q", "OK", got)
	}
	if requestCount != 2 {
		t.Fatalf("expected 2 requests (non-stream + stream fallback), got %d", requestCount)
	}
	if !streamSeen {
		t.Fatalf("expected stream fallback request to be sent")
	}
}

func TestCreateTextResponseConcatenatesResponsesStreamDeltasWithoutInjectingNewlines(t *testing.T) {
	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("expected valid json payload, got error: %v", err)
		}

		if payload["stream"] == true {
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"{\\\"intent\\\":\\\"host_lookup\\\",\"}\n\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"\\\"target_type\\\":\\\"ip\\\"}\"}\n\n"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_json_stream",
			"status": "completed",
			"error": null,
			"output": []
		}`))
	})

	got, err := client.CreateTextResponse(context.Background(), "Return JSON only.", "query host", 32)
	if err != nil {
		t.Fatalf("expected delta text to be concatenated into valid json string, got error: %v", err)
	}
	want := "{\"intent\":\"host_lookup\",\"target_type\":\"ip\"}"
	if got != want {
		t.Fatalf("expected concatenated stream json %q, got %q", want, got)
	}
}

func TestCreateTextResponseReturnsAsSoonAsResponsesStreamSignalsDone(t *testing.T) {
	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("expected valid json payload, got error: %v", err)
		}

		if payload["stream"] == true {
			w.Header().Set("Content-Type", "text/event-stream")
			flusher, _ := w.(http.Flusher)
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"OK\"}\n\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.done\",\"text\":\"OK\"}\n\n"))
			if flusher != nil {
				flusher.Flush()
			}
			time.Sleep(800 * time.Millisecond)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_stream_done",
			"status": "completed",
			"error": null,
			"output": []
		}`))
	})

	start := time.Now()
	got, err := client.CreateTextResponse(context.Background(), "Reply with OK only.", "health check", 8)
	if err != nil {
		t.Fatalf("expected stream fallback to finish on done event, got error: %v", err)
	}
	if got != "OK" {
		t.Fatalf("expected text %q, got %q", "OK", got)
	}
	if elapsed := time.Since(start); elapsed >= 700*time.Millisecond {
		t.Fatalf("expected stream reader to stop before server closed connection, took %s", elapsed)
	}
}

func TestProbePrefersStreamingHealthCheckForCompatibilityGateway(t *testing.T) {
	requestCount := 0
	streamSeen := false

	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("expected valid json payload, got error: %v", err)
		}

		if payload["stream"] == true {
			streamSeen = true
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"OK\"}\n\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.done\",\"text\":\"OK\"}\n\n"))
			return
		}

		time.Sleep(300 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_probe",
			"status": "completed",
			"error": null,
			"output": []
		}`))
	})
	client.httpClient.Timeout = 150 * time.Millisecond

	probe := client.Probe(context.Background())
	if !probe.Reachable {
		t.Fatalf("expected probe to succeed through streaming path, got last error %q", probe.LastError)
	}
	if probe.LastError != "" {
		t.Fatalf("expected no probe error, got %q", probe.LastError)
	}
	if requestCount != 1 {
		t.Fatalf("expected probe to use only one streaming request, got %d requests", requestCount)
	}
	if !streamSeen {
		t.Fatalf("expected streaming health check request to be used")
	}
}

func TestProbeStreamingHealthCheckOmitsReasoningEffortForLatency(t *testing.T) {
	client := newTestAIRuntimeClient(t, func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("expected valid json payload, got error: %v", err)
		}

		if payload["stream"] == true {
			if _, exists := payload["reasoning"]; exists {
				t.Fatalf("expected probe streaming request to omit reasoning config for lower latency")
			}
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"OK\"}\n\n"))
			_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.done\",\"text\":\"OK\"}\n\n"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"resp_probe_reasoning","status":"completed","error":null,"output":[]}`))
	})
	client.reasoningEffort = "xhigh"

	probe := client.Probe(context.Background())
	if !probe.Reachable {
		t.Fatalf("expected probe to succeed without reasoning, got last error %q", probe.LastError)
	}
}

func TestNormalizeAIRuntimeProbeErrorUsesFriendlyCopyForCompatibilityFailure(t *testing.T) {
	got := normalizeAIRuntimeProbeError(assertErrString("ai runtime completed successfully but did not return any text payload"))
	want := "模型网关已响应，但当前返回里没有标准文本内容，系统将先回退到内置逻辑。"
	if got != want {
		t.Fatalf("expected friendly compatibility message %q, got %q", want, got)
	}
}

func TestNormalizeAIRuntimeProbeErrorUsesFriendlyCopyForTimeout(t *testing.T) {
	got := normalizeAIRuntimeProbeError(assertErrString(`Post "https://beecode.cc/v1/responses": context deadline exceeded`))
	want := "模型网关探测超时，当前先回退到内置逻辑，请稍后重试。"
	if got != want {
		t.Fatalf("expected friendly timeout message %q, got %q", want, got)
	}
}

type errString string

func (e errString) Error() string {
	return string(e)
}

func assertErrString(message string) error {
	return errString(message)
}
