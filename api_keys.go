// Package plugindemo a demo plugin.
package traefik_plugin_api_keys

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// APIKeyValidatorConfig struct for plugin configuration
type APIKeyValidatorConfig struct {
	ValidAPIKeys       []string // List of valid API keys
	APIKeyHeader       string   // Header name to extract the API key
	UseAuthorization   bool     // Flag to use Authorization header for API key
	IgnorePaths        []string // Paths to ignore (no API key required)
    BlockPaths         []string // Paths to block (always deny access)
	RemoveHeader       bool     // Flag to remove the API key from the request header
}

// CreateConfig creates and initializes the configuration struct
func CreateConfig() *APIKeyValidatorConfig {
	return &APIKeyValidatorConfig{}
}

// APIKeyValidator struct representing the middleware
type APIKeyValidator struct {
	next           http.Handler
	name           string
	config         *APIKeyValidatorConfig
	ignorePathRegex []*regexp.Regexp
    blockPathRegex  []*regexp.Regexp
}

// New creates and returns a new APIKeyValidator instance
func New(ctx context.Context, next http.Handler, config *APIKeyValidatorConfig, name string) (http.Handler, error) {
	if config == nil {
		log.Println("APIKeyValidator configuration is nil")
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	if len(config.ValidAPIKeys) == 0 {
		log.Println("No valid API keys provided in configuration")
		return nil, fmt.Errorf("no valid API keys provided")
	}

	if config.APIKeyHeader == "" {
		config.APIKeyHeader = "X-API-Key" // Default header if not specified
	}

	ignorePathRegex, err := compilePathPatterns(config.IgnorePaths)
    if err != nil {
        return nil, fmt.Errorf("error compiling ignore paths: %v", err)
    }
    blockPathRegex, err := compilePathPatterns(config.BlockPaths)
    if err != nil {
        return nil, fmt.Errorf("error compiling block paths: %v", err)
    }

	return &APIKeyValidator{
		next:   next,
		name:   name,
		config: config,
		ignorePathRegex: ignorePathRegex,
        blockPathRegex:  blockPathRegex,
	}, nil

}

func compilePathPatterns(paths []string) ([]*regexp.Regexp, error) {
    var regexps []*regexp.Regexp
    for _, path := range paths {
        re, err := regexp.Compile(path)
        if err != nil {
			log.Printf("error compiling path pattern %s: %v", path, err)
			continue
        }
        regexps = append(regexps, re)
    }
    return regexps, nil
}

func (a *APIKeyValidator) isPathMatched(path string, patterns []*regexp.Regexp) bool {
    for _, pattern := range patterns {
        if pattern.MatchString(path) {
            return true
        }
    }
    return false
}


// log os.Stdout.WriteString("...") or os.Stderr.WriteString("...").

// ServeHTTP implements the HTTP handler interface
func (a *APIKeyValidator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("APIKeyValidator ServeHTTP"+ req.URL.Path+ " "+ req.Method+" "+ req.Header.Get("X-API-Key"))
	// Check if the path is in the ignore list
	isIgnorePath := a.isPathMatched(req.URL.Path, a.ignorePathRegex)

	if isIgnorePath {
		// Proceed with the next handler
		a.next.ServeHTTP(rw, req)
		return
	}

    // Check if the path is blocked
    if a.isPathMatched(req.URL.Path, a.blockPathRegex) {
        http.Error(rw, "Access to this path is blocked", http.StatusForbidden)
        return
    }

	apiKey := a.extractAPIKey(req)

	if apiKey == "" {
		http.Error(rw, "API key is missing or not provided correctly", http.StatusUnauthorized)
		fmt.Println("API key is missing or not provided correctly")
		return
	}

	if !a.isValidAPIKey(apiKey) {
		fmt.Println("Invalid API Key:" + apiKey)
		http.Error(rw, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	// Remove the API key from the request header
	if a.config.RemoveHeader {
		req.Header.Del(a.config.APIKeyHeader)
	}

		// Proceed with the next handler
	a.next.ServeHTTP(rw, req)
}

// extractAPIKey retrieves the API key from the request
func (a *APIKeyValidator) extractAPIKey(req *http.Request) string {
    if a.config.UseAuthorization {
        authHeader := req.Header.Get("Authorization")
        if strings.HasPrefix(authHeader, "Bearer ") {
            return strings.TrimPrefix(authHeader, "Bearer ")
        }
    } else {
        return req.Header.Get(a.config.APIKeyHeader)
    }
    return ""
}

// isValidAPIKey checks if the provided API key is valid
func (a *APIKeyValidator) isValidAPIKey(apiKey string) bool {
	for _, validKey := range a.config.ValidAPIKeys {
		if apiKey == validKey {
			return true
		}
	}
	return false
}