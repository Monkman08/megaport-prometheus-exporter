package megaport

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

const (
	prodTokenURL    = "https://auth-m2m.megaport.com/oauth2/token"
	stagingTokenURL = "https://oauth-m2m-staging.auth.ap-southeast-2.amazoncognito.com/oauth2/token"
	prodAPIURL      = "https://api.megaport.com"
	stagingAPIURL   = "https://api-staging.megaport.com"
)

// Client represents a Megaport API client
type Client struct {
	TokenURL   string
	APIURL     string
	APIKey     string
	APISecret  string
	HTTPClient *http.Client
	Token      string
	ExpiresAt  time.Time
}

// NewClient creates a new Megaport API client
func NewClient() (*Client, error) {
	apiKey := os.Getenv("MEGAPORT_API_KEY")
	apiSecret := os.Getenv("MEGAPORT_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		return nil, errors.New("MEGAPORT_API_KEY and MEGAPORT_API_SECRET must be set")
	}

	tokenURL := stagingTokenURL
	apiURL := stagingAPIURL
	if os.Getenv("MEGAPORT_ENV") == "prod" {
		tokenURL = prodTokenURL
		apiURL = prodAPIURL
	}

	return &Client{
		TokenURL:   tokenURL,
		APIURL:     apiURL,
		APIKey:     apiKey,
		APISecret:  apiSecret,
		HTTPClient: &http.Client{},
	}, nil
}

// GenerateToken generates a new access token
func (c *Client) GenerateToken() error {
	if time.Now().Before(c.ExpiresAt) {
		return nil
	}

	data := "grant_type=client_credentials"
	req, err := http.NewRequest("POST", c.TokenURL, bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.APIKey, c.APISecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to generate token")
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.Token = result.AccessToken
	c.ExpiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return nil
}

// DoRequest performs an authenticated request to the Megaport API
func (c *Client) DoRequest(method, endpoint string, body []byte) (*http.Response, error) {
	if err := c.GenerateToken(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.APIURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}
