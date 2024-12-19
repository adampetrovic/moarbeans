package woodroaster

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: "https://woodroaster.com.au",
	}
}

func (c *Client) RequestMagicLink(email string) error {
	data := url.Values{}
	data.Set("email", email)
	
	resp, err := c.httpClient.PostForm(c.baseURL+"/magic-link", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetNextOrderDate(sessionToken string) (time.Time, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/subscription/next-order", nil)
	if err != nil {
		return time.Time{}, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionToken))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	// Parse response and return next order date
	// Implementation depends on the actual response format
	return time.Now(), nil // placeholder
}

func (c *Client) SetNextOrderDate(sessionToken string, date time.Time) error {
	data := url.Values{}
	data.Set("next_order_date", date.Format("2006-01-02"))

	req, err := http.NewRequest("POST", c.baseURL+"/subscription/update-date", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", fmt.Sprintf("session=%s", sessionToken))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update order date: %d", resp.StatusCode)
	}

	return nil
} 