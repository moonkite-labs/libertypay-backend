package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RequestSigner handles HMAC-SHA256 request signing for additional security
type RequestSigner struct {
	secret []byte
}

// NewRequestSigner creates a new request signer with the given secret
func NewRequestSigner(secret string) *RequestSigner {
	return &RequestSigner{
		secret: []byte(secret),
	}
}

// SignRequest adds authentication signature to the request
func (rs *RequestSigner) SignRequest(req *http.Request, body []byte) error {
	if rs == nil {
		return nil // No signing configured
	}

	// Generate timestamp
	timestamp := time.Now().Unix()
	req.Header.Set("X-Timestamp", strconv.FormatInt(timestamp, 10))

	// Create string to sign
	stringToSign := rs.createStringToSign(req, body, timestamp)

	// Generate signature
	signature := rs.generateSignature(stringToSign)

	// Add signature header
	req.Header.Set("X-Signature", signature)

	return nil
}

// createStringToSign creates the canonical string to sign
func (rs *RequestSigner) createStringToSign(req *http.Request, body []byte, timestamp int64) string {
	var parts []string

	// HTTP method
	parts = append(parts, req.Method)

	// Path
	parts = append(parts, req.URL.Path)

	// Query parameters (sorted)
	if req.URL.RawQuery != "" {
		parts = append(parts, rs.canonicalizeQuery(req.URL.RawQuery))
	} else {
		parts = append(parts, "")
	}

	// Headers (specific headers only, sorted)
	parts = append(parts, rs.canonicalizeHeaders(req))

	// Body hash
	bodyHash := sha256.Sum256(body)
	parts = append(parts, hex.EncodeToString(bodyHash[:]))

	// Timestamp
	parts = append(parts, strconv.FormatInt(timestamp, 10))

	return strings.Join(parts, "\n")
}

// canonicalizeQuery sorts query parameters
func (rs *RequestSigner) canonicalizeQuery(query string) string {
	if query == "" {
		return ""
	}

	// Parse and sort query parameters
	params := strings.Split(query, "&")
	sort.Strings(params)
	return strings.Join(params, "&")
}

// canonicalizeHeaders includes specific headers in the signature
func (rs *RequestSigner) canonicalizeHeaders(req *http.Request) string {
	// Include specific headers that should be signed
	headersToSign := []string{
		"authorization",
		"content-type",
		"x-timestamp",
	}

	var headerParts []string

	for _, headerName := range headersToSign {
		value := req.Header.Get(headerName)
		if value != "" {
			headerParts = append(headerParts, fmt.Sprintf("%s:%s", headerName, strings.TrimSpace(value)))
		}
	}

	return strings.Join(headerParts, "\n")
}

// generateSignature generates HMAC-SHA256 signature
func (rs *RequestSigner) generateSignature(stringToSign string) string {
	h := hmac.New(sha256.New, rs.secret)
	h.Write([]byte(stringToSign))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature verifies a request signature (useful for testing)
func (rs *RequestSigner) VerifySignature(req *http.Request, body []byte, signature string) bool {
	if rs == nil {
		return true // No verification needed
	}

	timestampStr := req.Header.Get("X-Timestamp")
	if timestampStr == "" {
		return false
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	// Check timestamp is within acceptable range (5 minutes)
	now := time.Now().Unix()
	if abs(now-timestamp) > 300 {
		return false
	}

	// Generate expected signature
	stringToSign := rs.createStringToSign(req, body, timestamp)
	expectedSignature := rs.generateSignature(stringToSign)

	// Compare signatures (constant time comparison)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// signingTransport wraps an HTTP transport with request signing
type signingTransport struct {
	transport http.RoundTripper
	signer    *RequestSigner
}

func (st *signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read body if present
	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body for signing: %w", err)
		}
		// Restore body
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// Sign the request
	if err := st.signer.SignRequest(req, body); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return st.transport.RoundTrip(req)
}

// Helper function for absolute value
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}