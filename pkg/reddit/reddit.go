package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	authURL  = "https://www.reddit.com/api/v1/authorize"
	tokenURL = "https://www.reddit.com/api/v1/access_token"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	config     *oauth2.Config
	token      *oauth2.Token
	Statistics *Statistics
}

func NewClient(baseURL, clientID, clientSecret, redirectURI string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			Scopes:       []string{"read"},
		},
		Statistics: NewStatistics(),
	}
}

func (c *Client) AuthURL(responseType, duration, scope string) string {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("response_type", responseType)
	params.Set("state", generateState())
	params.Set("redirect_uri", "https://localhost/token")
	params.Set("duration", duration)
	params.Set("scope", scope)

	return authURL + "?" + params.Encode()
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	// body := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=https://localhost/token", code)
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "https://localhost/token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange authorization code: %s", resp.Status)
	}

	var token oauth2.Token
	if err := json.Unmarshal(res, &token); err != nil {
		return nil, fmt.Errorf("failed to decode response (%s): %w", string(res), err)
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("failed to decode response (%s)", string(res))
	}

	return &token, nil
}

// SetAccessToken sets the access token for the client
func (c *Client) SetAccessToken(token *oauth2.Token) {
	c.token = token
}

func (c *Client) GetPosts(ctx context.Context, subreddit string) ([]Post, error) {
	ts := c.config.TokenSource(ctx, c.token)
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	if token.Expiry.Before(time.Now()) {
		token, err = ts.Token()
		if err != nil {
			return nil, err
		}
		c.token = token
	}

	// Create authenticated HTTP client using refreshed token
	httpClient := oauth2.NewClient(ctx, ts)
	endpoint := fmt.Sprintf("/r/%s/hot?limit=2", subreddit)
	resp, err := httpClient.Get(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange authorization code: %s", resp.Status)
	}

	var response struct {
		Data struct {
			Children []struct {
				Data Post `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response (%s): %w", string(res), err)
	}

	var posts []Post
	for _, child := range response.Data.Children {
		posts = append(posts, child.Data)
		c.Statistics.Update(child.Data) // Update statistics for each post
	}

	return posts, nil
}

// Define structs for Reddit post data
type Post struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	UpVotes   int    `json:"ups"`
	DownVotes int    `json:"downs"`
}

// generateState generates a unique state string
func generateState() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	state := make([]byte, 32)
	for i := range state {
		state[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(state)
}
