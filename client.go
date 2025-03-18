package openroutergo

import (
	"net/http"
	"time"

	"github.com/eduardolat/openroutergo/internal/strutil"
)

const (
	defaultBaseURL = "https://openrouter.ai/api/v1"
	defaultTimeout = 3 * time.Minute
)

// Client represents a client for the OpenRouter API.
type Client struct {
	baseURL      string
	apiKey       string
	refererURL   string
	refererTitle string
	httpClient   *http.Client
}

// clientBuilder is a chainable builder for the OpenRouter client.
type clientBuilder struct {
	client *Client
}

// NewClient starts the creation of a new OpenRouter client.
func NewClient() *clientBuilder {
	return &clientBuilder{
		client: &Client{
			baseURL:      defaultBaseURL,
			apiKey:       "",
			refererURL:   "",
			refererTitle: "",
			httpClient:   &http.Client{Timeout: defaultTimeout},
		},
	}
}

// WithBaseURL sets a custom base URL for the API.
//
// If not set, the default base URL will be used: https://openrouter.ai/api/v1
func (b *clientBuilder) WithBaseURL(baseURL string) *clientBuilder {
	b.client.baseURL = strutil.RemoveTrailingSlashes(baseURL)
	return b
}

// WithAPIKey sets the API key for authentication.
func (b *clientBuilder) WithAPIKey(apiKey string) *clientBuilder {
	b.client.apiKey = apiKey
	return b
}

// WithRefererURL sets the referer URL for the API which identifies your app
// and allows it to be tracked and discoverable on OpenRouter.
//
// It uses the `HTTP-Referer` header.
//
//   - https://openrouter.ai/docs/api-reference/overview#headers
func (b *clientBuilder) WithRefererURL(refererURL string) *clientBuilder {
	b.client.refererURL = refererURL
	return b
}

// WithRefererTitle sets the referer title for the API which identifies your app
// and allows it to be discoverable on OpenRouter.
//
// It uses the `X-Title` header.
//
//   - https://openrouter.ai/docs/api-reference/overview#headers
func (b *clientBuilder) WithRefererTitle(refererTitle string) *clientBuilder {
	b.client.refererTitle = refererTitle
	return b
}

// WithHTTPClient sets a custom HTTP client for the API, this allows setting
// a custom timeout, proxy, etc.
//
// If not set, the default HTTP client will be used.
func (b *clientBuilder) WithHTTPClient(httpClient *http.Client) *clientBuilder {
	b.client.httpClient = httpClient
	return b
}

// WithTimeout sets a custom timeout for the HTTP client to be used for all requests.
//
// If not set, the default timeout of 3 minutes will be used.
func (b *clientBuilder) WithTimeout(timeout time.Duration) *clientBuilder {
	if b.client.httpClient == nil {
		b.client.httpClient = &http.Client{}
	}

	b.client.httpClient.Timeout = timeout
	return b
}

// Create builds and returns the OpenRouter client.
func (b *clientBuilder) Create() (*Client, error) {
	if b.client.baseURL == "" {
		return nil, ErrBaseURLRequired
	}

	if b.client.apiKey == "" {
		return nil, ErrAPIKeyRequired
	}

	return b.client, nil
}
