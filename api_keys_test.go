package traefik_plugin_api_kyes_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	apikeys "github.com/a8851625/traefik-plugin-api-kyes"
)

// TestAPIKeyValidator tests various scenarios for the APIKeyValidator plugin

func TestAPIKeyValidator(t *testing.T) {
    cfg := apikeys.CreateConfig()
    cfg.IgnorePaths = []string{"/ignore"}
    cfg.BlockPaths = []string{"/block.*"}
    cfg.ValidAPIKeys = []string{"sk-valid-key-sample"}
    cfg.APIKeyHeader = "X-API-Key"
	cfg.UseAuthorization = false
    cfg.RemoveHeader = true

    ctx := context.Background()
    next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

    handler, err := apikeys.New(ctx, next, cfg, "api-key-validator")
    if err != nil {
        t.Fatal(err)
    }

    // Test scenarios
    testScenarios := []struct {
        description string
        path        string
        apiKey      string
        expectedStatus int
    }{
        {"Ignore Path", "/ignore", "", http.StatusOK},
        {"Block Path", "/block1", "", http.StatusForbidden},
        {"Valid API Key", "/normal", "sk-valid-key-sample", http.StatusOK},
        {"Invalid API Key", "/normal", "invalid-key", http.StatusUnauthorized},
    }

    for _, scenario := range testScenarios {
        t.Run(scenario.description, func(t *testing.T) {
            recorder := httptest.NewRecorder()

            req, err := http.NewRequestWithContext(ctx, http.MethodGet, scenario.path, nil)
            if err != nil {
                t.Fatal(err)
            }

            // Set API key if provided
            if scenario.apiKey != "" {
                req.Header.Set(cfg.APIKeyHeader, scenario.apiKey)
            }

            handler.ServeHTTP(recorder, req)

            if recorder.Result().StatusCode != scenario.expectedStatus {
                t.Errorf("unexpected status for %s: got %v want %v",
                    scenario.description, recorder.Result().StatusCode, scenario.expectedStatus)
            }

            // Check header removal
            if _, exists := req.Header[cfg.APIKeyHeader]; exists && scenario.expectedStatus == http.StatusOK {
                t.Errorf("header %s was not removed as expected", cfg.APIKeyHeader)
            }
        })
    }
}
