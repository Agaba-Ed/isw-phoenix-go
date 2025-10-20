package isw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"isw-phoenix-go/internal/logger"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient handles HTTP requests
type HTTPClient struct {
	client     *http.Client
	config     *Configuration
	privateKey string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(config *Configuration, privateKey string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		config:     config,
		privateKey: privateKey,
	}
}

// Get performs GET request
func (h *HTTPClient) Get(endpoint string, params map[string]string, authToken, sessionKey string) (*APIResponse, error) {
	fullURL := h.config.APIURL + endpoint

	// Add query parameters
	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	headers := h.buildHeaders("GET", fullURL, authToken, sessionKey, nil)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return h.handleResponse(resp)
}

// GetV2 performs GET request to billers API
func (h *HTTPClient) GetV2(endpoint string, params map[string]string, authToken, sessionKey string) (*APIResponse, error) {
	fullURL := h.config.BillersAPIURL + endpoint

	// Add query parameters
	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	headers := h.buildHeaders("GET", fullURL, authToken, sessionKey, nil)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return h.handleResponse(resp)
}

// Post performs POST request
func (h *HTTPClient) Post(endpoint string, data interface{}, authToken, sessionKey string, additionalParams *AdditionalParameters) (*APIResponse, error) {
	fullURL := h.config.APIURL + endpoint

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	headers := h.buildHeaders("POST", fullURL, authToken, sessionKey, additionalParams)

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return h.handleResponse(resp)
}

// buildHeaders builds HTTP headers for requests
func (h *HTTPClient) buildHeaders(method, fullURL, authToken, sessionKey string, additionalParams *AdditionalParameters) map[string]string {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	nonce := generateNonce()

	encryptedAuthToken, _ := EncryptAuthToken(authToken, sessionKey)

	signature, _ := GenerateSignatureHeader(method, fullURL, timestamp, nonce, h.config.ClientID, h.config.ClientSecretKey, h.privateKey, additionalParams)

	return map[string]string{
		"Authorization": GenerateAuthorizationHeader(h.config.ClientID),
		"Signature":     signature,
		"Content-Type":  "application/json",
		"Nonce":         nonce,
		"Timestamp":     timestamp,
		"AuthToken":     encryptedAuthToken,
	}
}

// handleResponse handles HTTP response
func (h *HTTPClient) handleResponse(resp *http.Response) (*APIResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// Fallback: handle cases where responseCode may be a number instead of a string
		var generic map[string]interface{}
		if err2 := json.Unmarshal(body, &generic); err2 != nil {
			// Log the non-JSON response for debugging (status and a trimmed body)
			trim := body
			if len(trim) > 512 {
				trim = trim[:512]
			}
			logger.Log.WithFields(map[string]interface{}{
				"status":      resp.Status,
				"contentType": resp.Header.Get("Content-Type"),
				"bodySnippet": string(trim),
			}).Warn("Upstream returned non-JSON response")
			return &APIResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to parse response: %v", err),
			}, nil
		}

		// Normalize fields into APIResponse
		apiResp = APIResponse{}

		// responseCode -> string
		if rc, ok := generic["responseCode"]; ok {
			switch v := rc.(type) {
			case string:
				apiResp.ResponseCode = v
			case float64:
				// format without decimals
				apiResp.ResponseCode = fmt.Sprintf("%.0f", v)
			default:
				apiResp.ResponseCode = fmt.Sprintf("%v", v)
			}
		}

		if rm, ok := generic["responseMessage"].(string); ok {
			apiResp.ResponseMessage = rm
		}
		if rr, ok := generic["requestReference"]; ok {
			// could be null or string
			if s, ok := rr.(string); ok {
				apiResp.RequestReference = s
			}
		}
		if r, ok := generic["response"]; ok {
			apiResp.Response = r
		}

		// infer success when code is 90000
		apiResp.Success = (apiResp.ResponseCode == "90000")
	}

	if resp.StatusCode >= 400 {
		apiResp.Success = false
		if apiResp.Error == "" {
			apiResp.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
	}

	return &apiResp, nil
}

// generateNonce generates a random nonce
func generateNonce() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
